package database

import (
	"errors"
	"database/sql"
	"github.com/google/uuid"
)

//Player - business object to represent database la_player record
type Player struct{
	PlayerID			uuid.UUID 	
	PlayerName			string	
	PlayerBio	 		string	
	PlayerProfileLink 	string	
	PlayerType			int
	PlayerPicture		[]byte
	IsActive			bool
}

//UserPlayerMap - business object to represent database la_user_player_map record
type UserPlayerMap struct{
	MapID				uuid.UUID
	UserID				uuid.UUID
	PlayerID			uuid.UUID
}

//PlayerStore - Player db store contract
type PlayerStore interface{
	GetPlayerByUserUUID(userUUID uuid.UUID) (*Player, error)
	GetPlayerByPlayerUUID(playerUUID uuid.UUID) (*Player, error)
	UpdatePlayerInfoForUser(player *Player, userUUID uuid.UUID) error
}

//GetPlayerDBStore - get player database store for tests
func GetPlayerDBStore(_db *sql.DB) PlayerStore{
	playerstore := new(playerStoreDbImpl)
	playerstore.db = _db
	return playerstore
}

type playerStoreDbImpl struct{
	db *sql.DB
}

const (
	selectUserPlayerMappingByUserUUIDQuery = "SELECT player_id FROM la_schema.la_user_player_map WHERE user_id = $1"
	fetchPlayerInfoByPlayerUUIDQuery = "SELECT player_name, player_bio, player_profile_link, player_type, player_photo, is_active FROM la_schema.la_player WHERE player_id = $1"
	upsertPlayerInfo = `INSERT INTO la_schema.la_player (player_id, player_name, player_bio, player_profile_link, player_type, player_photo, is_active) 
						VALUES($1, $2, $3, $4, $5, $6, $7)
						ON CONFLICT (player_id)
						DO UPDATE
						SET player_name = EXCLUDED.player_name,
						player_bio = EXCLUDED.player_bio,
						player_profile_link = EXCLUDED.player_profile_link,
						player_type = EXCLUDED.player_type,
						player_photo = EXCLUDED.player_photo,
						is_active = EXCLUDED.is_active
						RETURNING player_id`
	upsertUserPlayerMapping = `INSERT INTO la_schema.la_user_player_map(user_id, player_id)
								VALUES($1, $2)
								ON CONFLICT(user_id) DO NOTHING`
)

func (ps *playerStoreDbImpl)GetPlayerByUserUUID(userUUID uuid.UUID) (*Player, error) {
	if ps.db == nil {
		return nil, errors.New("database object con not be nil")
	}
	up := UserPlayerMap{UserID : userUUID}
	err := ps.db.QueryRow(selectUserPlayerMappingByUserUUIDQuery,userUUID).Scan(&up.PlayerID)
	if err != nil{
		return nil, err
	}
	player := new(Player)
	player.PlayerID = up.PlayerID
	err = ps.db.QueryRow(fetchPlayerInfoByPlayerUUIDQuery,up.PlayerID).Scan(&player.PlayerName, &player.PlayerBio, &player.PlayerProfileLink, &player.PlayerType, &player.PlayerPicture, &player.IsActive)
	if err != nil{
		return nil, err
	}
	return player, nil
}


func (ps *playerStoreDbImpl)GetPlayerByPlayerUUID(playerUUID uuid.UUID) (*Player, error) {
	if ps.db == nil {
		return nil, errors.New("database object con not be nil")
	}
	player := new(Player)
	player.PlayerID = playerUUID
	err := ps.db.QueryRow(fetchPlayerInfoByPlayerUUIDQuery,playerUUID).Scan(&player.PlayerName, &player.PlayerBio, &player.PlayerProfileLink, &player.PlayerType, &player.PlayerPicture, &player.IsActive)
	if err != nil{
		return nil, err
	}
	return player, nil
}


func (ps *playerStoreDbImpl)UpdatePlayerInfoForUser(player *Player, userUUID uuid.UUID) error{
	if ps.db == nil {
		return errors.New("database object can not be nil")
	}
	if player == nil {
		return errors.New("player object can not be nil")
	}
	_, err := ps.GetPlayerByUserUUID(userUUID)
	if err != nil && err == sql.ErrNoRows{
		player.PlayerID = uuid.New()
	} else if err != nil{
		return err
	}
	_, err = ps.db.Exec(upsertPlayerInfo, player.PlayerID, player.PlayerName, player.PlayerBio, 
						player.PlayerProfileLink, player.PlayerType,
						player.PlayerPicture, player.IsActive)
	if err != nil{
		return err
	}
	_, err = ps.db.Exec(upsertUserPlayerMapping, userUUID, player.PlayerID)
	if err != nil{
		return err
	}
	return nil
}



//GetPlayerMockStore - get player mock store for tests
func GetPlayerMockStore() PlayerStore{
	playerstore := new(playerStoreMockImpl)
	playerstore.playerIDInfoMap = make(map[uuid.UUID]*Player)
	playerstore.userPlayerIDMap = make(map[uuid.UUID]uuid.UUID)
	return playerstore
}

type playerStoreMockImpl struct{
	playerIDInfoMap	map[uuid.UUID]*Player
	userPlayerIDMap map[uuid.UUID]uuid.UUID
}

func (ps *playerStoreMockImpl)GetPlayerByUserUUID(userUUID uuid.UUID) (*Player, error) {
	if playerUUID, found := ps.userPlayerIDMap[userUUID]; found == true{
		if player, playerfound := ps.playerIDInfoMap[playerUUID]; playerfound == true{
			return player, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (ps *playerStoreMockImpl)GetPlayerByPlayerUUID(playerUUID uuid.UUID) (*Player, error) {
	if player, playerfound := ps.playerIDInfoMap[playerUUID]; playerfound == true{
		return player, nil
	}
	return nil, sql.ErrNoRows
}

func (ps *playerStoreMockImpl)UpdatePlayerInfoForUser(player *Player, userUUID uuid.UUID) error{
	if playerUUID, found := ps.userPlayerIDMap[userUUID]; found == true{
		player.PlayerID = ps.playerIDInfoMap[playerUUID].PlayerID	// preserve playerid
		ps.playerIDInfoMap[playerUUID] = player
		return nil
	}
	player.PlayerID = uuid.New()
	ps.userPlayerIDMap[userUUID] = player.PlayerID
	ps.playerIDInfoMap[player.PlayerID] = player
	return nil
}