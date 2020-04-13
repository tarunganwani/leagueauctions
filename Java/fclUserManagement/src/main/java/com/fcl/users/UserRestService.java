package com.fcl.users;

import com.fasterxml.jackson.databind.ObjectMapper;

import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.StreamingOutput;


@Path("/users")
public class UserRestService extends JsonMap<UserAccount> {
    private UserServices userServices;

    public UserRestService() {
        super();
        this.userServices = new UserServicesImpl();
    }

    @POST
    @Path("/login")
    @Produces(MediaType.APPLICATION_JSON)
    @Consumes(MediaType.APPLICATION_JSON)
    public Response loginUser(String loginUser) {
        System.out.println("User : " + loginUser);
        UserAccount userAccount = mapToJson(loginUser,UserAccount.class);
        System.out.println("Email ID = " + userAccount.getEmailId());
        boolean isAvail = userServices.checkUser(userAccount.getEmailId());
        String respJSON = "{\"user\":\"" + isAvail + "\"}";
        return Response.ok(respJSON).build();
    }
}
