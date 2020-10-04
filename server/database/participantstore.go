package database

import (
	// "log"
	// "fmt"
	// "time"
	"errors"
	"database/sql"
	// "database/sql/driver"
	"github.com/google/uuid"
	// "github.com/lib/pq"
)


const (
	createParticipantQuery 			= `INSERT INTO la_schema.la_participant(participant_id, player_id, auction_id, participant_role, category_id) values($1, $2, $3, $4, $5)`
	updateParticipantQuery 			= `UPDATE la_schema.la_participant SET participant_role = $1, category_id = $2 WHERE participant_id = $3 and auction_id = $4`
	fetchAllparticipantsQuery 		= `SELECT participant_id, player_id, participant_role, category_id FROM la_schema.la_participant WHERE auction_id = $1`
)

//ParticipantRoleType - type to define participant role
type ParticipantRoleType int

const (
	//ViewerRole - viewer only role
	ViewerRole ParticipantRoleType = iota
	//PlayerRole - player under the hammer
	PlayerRole
	//CaptainRole - captain or bidder
	CaptainRole
	//AuctioneerRole - auction controller
	AuctioneerRole
	//AuctioneerPlayerRole - auction controller and player under the hammer 
	AuctioneerPlayerRole
)

//Participant - participant databas object
type Participant struct{
	ParticipantUUID uuid.UUID
	PlayerUUID uuid.UUID
	AuctionBoardUUID uuid.UUID
	ParticipantRole ParticipantRoleType 
	CategoryUUID uuid.UUID
}


//ParticipantStore - participant store interface
type ParticipantStore interface{
	CreateParticipant(p *Participant) error	
	UpdateParticipant(p *Participant) error	
	FetchAllParticipants(auctionUUID uuid.UUID) ([]*Participant, error)
}

//GetParticipantDBStore - get participant database store
func GetParticipantDBStore(_db *sql.DB) ParticipantStore{
	participantstore := new(ParticipantDBStore)
	participantstore.db = _db
	return participantstore
}

//ParticipantDBStore - particpant database store
type ParticipantDBStore struct{
	db *sql.DB
}

//CreateParticipant - insert participant in database
func (ps *ParticipantDBStore)CreateParticipant(p *Participant) error{
	if ps.db == nil {
		return errors.New("database object can not be nil")
	}
	if p == nil {
		return errors.New("participant object can not be nil")
	}
	p.ParticipantUUID = uuid.New()
	_, err := ps.db.Exec(createParticipantQuery, p.ParticipantUUID, p.PlayerUUID,
							p.AuctionBoardUUID, ViewerRole, uuid.Nil)
	return err
}

//UpdateParticipant - update particpant attributes(role for now)
func (ps *ParticipantDBStore)UpdateParticipant(p *Participant) error{
	if ps.db == nil {
		return errors.New("database object can not be nil")
	}
	if p == nil {
		return errors.New("participant object can not be nil")
	}
	_, err := ps.db.Exec(updateParticipantQuery, p.ParticipantRole, p.CategoryUUID, 
												p.ParticipantUUID, p.AuctionBoardUUID)
	return err
}

//FetchAllParticipants - get participant list
func (ps *ParticipantDBStore)FetchAllParticipants(auctionUUID uuid.UUID) ([]*Participant, error){
	if ps.db == nil {
		return nil, errors.New("database object can not be nil")
	}
	if auctionUUID == uuid.Nil{
		return nil, errors.New("auction uuid can not be nil")
	}
	participantRows, err := ps.db.Query(fetchAllparticipantsQuery,auctionUUID)
	if err != nil{
		return nil, err
	}
	defer participantRows.Close()
	outParticipantSlice := make([]*Participant, 0)
	for participantRows.Next() {
		participant := new(Participant)
		participant.AuctionBoardUUID = auctionUUID
		if err := participantRows.Scan(&participant.ParticipantUUID, &participant.PlayerUUID, 
										&participant.ParticipantRole, &participant.CategoryUUID); err != nil {
			return nil, err
		}
		outParticipantSlice = append(outParticipantSlice , participant)
	}
	return outParticipantSlice, nil
}


//GetParticipantMockStore - get participant mock store for tests
func GetParticipantMockStore(inMockAuctionstore AuctionStore) ParticipantStore{
	participantstore := new(ParticipantMockStore)
	participantstore.auctionIDToParticipantDictMap = make(map[uuid.UUID]map[uuid.UUID]*Participant)
	participantstore.mockAuctionStore = inMockAuctionstore
	return participantstore
}

//ParticipantMockStore - particpant mock store
type ParticipantMockStore struct{
	auctionIDToParticipantDictMap map[uuid.UUID]map[uuid.UUID]*Participant
	mockAuctionStore AuctionStore
}

//CreateParticipant - insert participant in database
func (ps *ParticipantMockStore)CreateParticipant(p *Participant) error{
	if p.AuctionBoardUUID == uuid.Nil || ps.mockAuctionStore == nil{
		return errors.New("ParticipantMockStore - auction UUID or store can not be nil")
	}
	if p == nil {
		return errors.New("ParticipantMockStore - participantcan not be nil")
	}
	if _, err := ps.mockAuctionStore.GetAuctionBoardInfo(p.AuctionBoardUUID); err != nil{
		return err
	}
	ps.auctionIDToParticipantDictMap[p.AuctionBoardUUID] = make(map[uuid.UUID]*Participant)
	if  participantDict, listFound := ps.auctionIDToParticipantDictMap[p.AuctionBoardUUID]; listFound == true{
		if _, participantFound := participantDict[p.ParticipantUUID]; participantFound == false{
			p.ParticipantUUID = uuid.New()
			ps.auctionIDToParticipantDictMap[p.AuctionBoardUUID][p.ParticipantUUID] = p
			return nil
		}
		return errors.New("ParticipantMockStore - particpant already exists")
	}
	return sql.ErrNoRows
}

//UpdateParticipant - update particpant attributes(role for now)
func (ps *ParticipantMockStore)UpdateParticipant(p *Participant) error{
	if p.AuctionBoardUUID == uuid.Nil{
		return errors.New("ParticipantMockStore - auction UUID can not be nil")
	}
	if p == nil || p.PlayerUUID == uuid.Nil{
		return errors.New("ParticipantMockStore - participant or participant uuid can not be nil")
	}
	if  participantDict, listFound := ps.auctionIDToParticipantDictMap[p.AuctionBoardUUID]; listFound == true{
		if _, participantFound := participantDict[p.ParticipantUUID]; participantFound == true{
			ps.auctionIDToParticipantDictMap[p.AuctionBoardUUID][p.ParticipantUUID] = p
			return nil
		}
		return errors.New("ParticipantMockStore - particpant does not exist")
	}
	return sql.ErrNoRows
	//return errors.New("ParticipantMockStore - invalid auction uuid ")
}

//FetchAllParticipants - get participant list
func (ps *ParticipantMockStore)FetchAllParticipants(auctionUUID uuid.UUID) ([]*Participant, error){
	if auctionUUID == uuid.Nil{
		return nil, errors.New("ParticipantMockStore - auction UUID can not be nil")
	}
	i := 0
	if  participantDict, listFound := ps.auctionIDToParticipantDictMap[auctionUUID]; listFound == true{
		participantList := make([]*Participant, len(ps.auctionIDToParticipantDictMap))
		for _, participant := range participantDict{
			participantList[i] = participant
		}
		return participantList, nil
	}
	return nil, sql.ErrNoRows
	// return nil, errors.New("ParticipantMockStore - auction uuid does not exist")
}