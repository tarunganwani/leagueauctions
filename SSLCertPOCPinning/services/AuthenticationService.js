import { BehaviorSubject } from 'rxjs';

//import config from 'config';
import axios from 'axios';

//const currentUserSubject = new BehaviorSubject(JSON.parse(localStorage.getItem('currentUser')));

export const authenticationService = {
    login,
    logout,
    //currentUser: currentUserSubject.asObservable(),
    //get currentUserValue() { return currentUserSubject.value }
};

function login(username, password) {
    const requestOptions = {
        "user_id": username,
        "user_password": password
    };

    return axios.post(`https://localhost:8081/user/login`, requestOptions)
        .then(function (response) {
            console.log(response);
            // if (response.status === 200) {
            //     return response.data;
            // } else {
            //     console.log("Some error ocurred " + response.data);
            // }
        })
        .catch(function (error) {
            console.log(error);
        })
        .then(tokens => {
            // console.log("Users : " + JSON.stringify(tokens))
            // // store user details and jwt token in local storage to keep user logged in between page refreshes
            // localStorage.setItem("tokens", JSON.stringify(tokens));
            // currentUserSubject.next(tokens);
            // return tokens;
        });
}

function logout() {
    // remove user from local storage to log user out
    localStorage.removeItem('currentUser');
    currentUserSubject.next(null);
}
