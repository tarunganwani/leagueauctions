import React from 'react';
import TextField from "@material-ui/core/TextField";

const OtpScreen = () => {
    return (
        <div>
            <p> Otp has been send to your Email Address </p>
            <form noValidate autoComplete="off">
                <TextField id="otp" label="OTP" color="secondary" />
            </form>
        </div>
    )
}

export default OtpScreen