package com.fcl.users;

import com.fasterxml.jackson.annotation.JsonProperty;

public class UserAccount {
    @JsonProperty
    String id;
    @JsonProperty
    String emailId;
    @JsonProperty
    String password;
    @JsonProperty
    String password_salt;
    @JsonProperty
    String isActive;

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getPassword_salt() {
        return password_salt;
    }

    public void setPassword_salt(String password_salt) {
        this.password_salt = password_salt;
    }

    public String getEmailId() {
        return emailId;
    }

    public void setEmailId(String emailId) {
        this.emailId = emailId;
    }
}
