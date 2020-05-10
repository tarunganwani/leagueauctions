import React from 'react';
import Login from '../Login/Login';
import { Switch, Route } from 'react-router-dom';
import Register from '../Register/Register';

export default function Loginscreen() {
    return (
        <div>
            <Switch>
                <Route path="/login" exact component={Login} />
                <Route path="/register" exact component={Register} />
            </Switch>
        </div>
    );
}