export const SIGNUP = 'SIGNUP';
export const LOGIN = 'LOGIN';

export const signup = (email, password) => {
    return async dispatch => {
        const response = await fetch(
            'http://localhost:8080/user/register',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    user_id: email,
                    user_password: password,
                })
            }
        );

        if (response.status != 200) {
            throw new Error('Something went wrong!');
        }

        const resData = await response.json();
        console.log(resData);
        dispatch({ type: SIGNUP });
    };
};

export const login = (email, password) => {
    return async dispatch => {
        const response = await fetch(
            'http:192.168.1.22:8080/user/login',
            {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    user_id: email,
                    user_password: password,
                    //returnSecureToken: true
                })
            }
        );

        if (response.status === 200) {
            const resData = await response.json();
            console.log(resData);
            dispatch({ type: LOGIN , token: resData.login_token , userId: email});
        } else {
            const errorResData = await response.json();
            const errorId = errorResData.error.message;
            let message = 'Something went wrong!';
            if (errorId === 'EMAIL_NOT_FOUND') {
              message = 'This email could not be found!';
            } else if (errorId === 'INVALID_PASSWORD') {
              message = 'This password is not valid!';
            }
            throw new Error(message);
        }
    };
};