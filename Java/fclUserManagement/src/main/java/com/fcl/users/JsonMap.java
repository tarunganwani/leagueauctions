package com.fcl.users;

import com.fasterxml.jackson.databind.ObjectMapper;

public class JsonMap<T> {
    private final ObjectMapper objectMapper;

    JsonMap() {
        objectMapper = new ObjectMapper();
    }

    public T mapToJson(String user, Class<T> clazz) {
        try {
            return objectMapper.readValue(user, clazz);
        } catch (Exception e) {
            e.printStackTrace();
            throw new IllegalArgumentException("Invalid JSON passed" + user);
        }
    }
}
