import React, { useState } from "react";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import { makeStyles } from "@material-ui/core/styles";
import { Route, Link, Switch } from 'react-router-dom';
import axios from "axios";
import Register from '../Register/Register';

const useStyles = makeStyles((theme) => ({
    root: {
        "& > *": {
            margin: theme.spacing(1),
            width: "25ch",
        },
    },
}));

export default function Login() {
    const classes = useStyles();
    const [state, setState] = useState({
        email: "",
        password: "",
    });

    const handleChange = (e) => {
        const { id, value } = e.target;
        setState((prevState) => ({
            ...prevState,
            [id]: value,
        }));
    };

    const handleSubmitClick = (event) => {
        if (state.email.length && state.password.length) {
            //props.showError(null);
            //https://localhost:4000/myapp/api/login
            console.log(state.email + "   " + state.password);
            // const payload = {
            //     "email": state.email,
            //     "password": state.password,
            // }
            axios
                .get("http://localhost:8080/")
                .then(function (response) {
                    if (response.data.code === 200) {
                        console.log("200");
                    } else {
                        console.log("Some error ocurred " + response.data);
                    }
                })
                .catch(function (error) {
                    console.log(error);
                });
        } else {
            console.log("Please enter valid username and password");
        }
    };

    return (
        <form className={classes.root} noValidate autoComplete="off" >
            <TextField id="email"
                label="Username"
                onChange={handleChange}
                value={state.email} />
            <br />
            <TextField id="password"
                label="Password"
                type="password"
                onChange={handleChange}
                value={state.password}
                autoComplete="current-password" />
            <br />
            <Button
                variant="contained"
                color="primary"
                onClick={handleSubmitClick}> Submit </Button>
            <br />
            Not registered yet, <Link to="/register"> Register Now </Link>
            <Switch>
                <Route path="/register" exact component={Register} />
            </Switch>
        </form>

    );
}