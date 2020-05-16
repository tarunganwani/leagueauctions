import React from 'react';
import TextField from "@material-ui/core/TextField";

const RegistrationPage = () => {
    return (
        <div>
            <form noValidate autoComplete="off" >
                <TextField id="email"
                    label="Username"
                //onChange={handleChange}
                //value={state.email}
                />
                <br />
                <TextField id="password"
                    label="Password"
                    type="password"
                    //onChange={handleChange}
                    //value={state.password}
                    autoComplete="current-password" />
                <br />
            </form>
        </div>
    )
}

export default RegistrationPage