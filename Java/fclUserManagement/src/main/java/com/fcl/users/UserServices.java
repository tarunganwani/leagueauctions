package com.fcl.users;

public interface UserServices {

    public boolean closeConnection();

    public void updateUser(String id, PlayerProfile playerProfile);

    public void createUser(PlayerProfile playerProfile);

    public void deleteUser();

    public boolean checkUser(String emailId);

}
