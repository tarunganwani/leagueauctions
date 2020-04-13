package com.fcl.users;

import com.fasterxml.jackson.annotation.JsonProperty;

public class PlayerProfile {
    @JsonProperty
    String id;
    @JsonProperty
    String emailId;
    @JsonProperty
    String password;
    @JsonProperty
    String cricHeroesProfile;
    @JsonProperty
    String isActive;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getIsActive() {
        return isActive;
    }

    public void setIsActive(String isActive) {
        this.isActive = isActive;
    }

    public String getEmailId() {
        return emailId;
    }

    public void setEmailId(String emailId) {
        this.emailId = emailId;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public String getCricHeroesProfile() {
        return cricHeroesProfile;
    }

    public void setCricHeroesProfile(String cricHeroesProfile) {
        this.cricHeroesProfile = cricHeroesProfile;
    }
}
