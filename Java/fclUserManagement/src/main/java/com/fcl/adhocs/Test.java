package com.fcl.adhocs;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fcl.users.UserAccount;

import java.io.IOException;

public class Test {
    ObjectMapper objectMapper;

    Test() throws JsonProcessingException {
        objectMapper = new ObjectMapper();
    }

     void process() throws IOException {
        String json = "{ \"emailId\" : \"Black\", \"mobileNo\" : \"BMW\" }";
        UserAccount car = objectMapper.readValue(json, UserAccount.class);
        System.out.println(car.getEmailId() + "\t");
    }
}
