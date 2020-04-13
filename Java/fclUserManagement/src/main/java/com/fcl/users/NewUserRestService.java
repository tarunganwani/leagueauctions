package com.fcl.users;

import javax.ws.rs.Consumes;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;

@Path("/new")
public class NewUserRestService extends JsonMap<PlayerProfile> {
    private UserServices userServices;

    public NewUserRestService() {
        super();
        this.userServices = new UserServicesImpl();
    }

    @POST
    @Path("/register")
    @Produces(MediaType.APPLICATION_JSON)
    @Consumes(MediaType.APPLICATION_JSON)
    public Response registerUser(String newUser) {
        System.out.println("User For Registration :" + newUser);
        PlayerProfile playerProfile = mapToJson(newUser, PlayerProfile.class);
        System.out.println("Email ID :" + playerProfile.getEmailId());
        //boolean isAvail = userServices.checkUser(playerProfile.getEmailId());
        //if(!isAvail) {
        userServices.createUser(playerProfile);
        //} else {
        //    System.out.println("user already exists please login");
        //}
        return Response.ok().build();
    }

    @POST
    @Path("/update")
    @Produces(MediaType.APPLICATION_JSON)
    @Consumes(MediaType.APPLICATION_JSON)
    public Response updateUser(String user) {
        System.out.println("User For Updation :" + user);
        PlayerProfile playerProfile = mapToJson(user, PlayerProfile.class);
        System.out.println("Email ID :" + playerProfile.getEmailId());
        //boolean isAvail = userServices.checkUser(playerProfile.getEmailId());
        //if(!isAvail) {
        userServices.updateUser(playerProfile.getId(), playerProfile);
        //} else {
        //    System.out.println("user already exists please login");
        //}
        return Response.ok().build();
    }

}
