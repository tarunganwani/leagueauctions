package com.fcl.users;

import java.sql.*;

public class UserServicesImpl implements UserServices {
    private Connection connection = null;

    public UserServicesImpl() {
    }

    public void connectDb() throws ClassNotFoundException {
        Class.forName("org.sqlite.JDBC");

        try {
            // create a database connection
            connection = DriverManager.getConnection("jdbc:sqlite:F:/Database/SQLite DB/mydb.db");
        } catch (SQLException e) {
            // if the error message is "out of memory",
            // it probably means no database file is found
            System.err.println(e.getMessage());
        }
        /*finally {

        }*/
    }

    public boolean checkConnection() {
        return connection != null;
    }

    @Override
    public boolean closeConnection() {
        try {
            if (checkConnection())
                connection.close();
            return true;
        } catch (SQLException e) {
            // connection close failed.
            e.printStackTrace();
        }
        return false;
    }

    @Override
    public void updateUser(String id, PlayerProfile playerProfile) {
        if (!checkConnection()) {
            try {
                connectDb();
            } catch (ClassNotFoundException e) {
                e.printStackTrace();
            }
        }
        PreparedStatement statement = null;
        try {
            statement = connection.prepareStatement("update UserProfile set emailId = ? , cricHeroesProfile = ? where id = ?");
            statement.setString(1, playerProfile.getEmailId());
            statement.setString(2, playerProfile.getCricHeroesProfile());
            statement.setString(3, id);
            statement.setQueryTimeout(30);
            int i = statement.executeUpdate();
            System.out.println(i + " records updated");
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void createUser(PlayerProfile playerProfile) {
        if (!checkConnection()) {
            try {
                connectDb();
            } catch (ClassNotFoundException e) {
                e.printStackTrace();
            }
        }
        PreparedStatement statement = null;
        try {
            statement = connection.prepareStatement("insert into UserProfile values(?,?,?,?,?)");
            statement.setString(1, playerProfile.getId());
            statement.setString(2, playerProfile.getEmailId());
            statement.setString(3, playerProfile.getPassword());
            statement.setString(4, playerProfile.getCricHeroesProfile());
            statement.setString(5, playerProfile.isActive);
            statement.setQueryTimeout(30);
            int i = statement.executeUpdate();
            System.out.println(i + " records inserted");

        } catch (SQLException e) {
            e.printStackTrace();
        }
    }

    @Override
    public void deleteUser() {

    }

    @Override
    public boolean checkUser(String emailId) {
        boolean usernameExists = false;
        if (!checkConnection()) {
            try {
                connectDb();
            } catch (ClassNotFoundException e) {
                e.printStackTrace();
            }
        }

        PreparedStatement statement = null;
        try {
            statement = connection.prepareStatement("select * from person where name = ?");
            statement.setString(1, emailId);
            statement.setQueryTimeout(30);  // set timeout to 30 sec.

            ResultSet rs = statement.executeQuery();
            if (rs.next()) {
                usernameExists = true;
            }

        } catch (SQLException e) {
            e.printStackTrace();
        }
        return usernameExists;
    }
}
