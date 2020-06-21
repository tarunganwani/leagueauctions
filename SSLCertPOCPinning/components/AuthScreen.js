import React, { useState } from 'react';
import {
    TextInput,
    StyleSheet,
    KeyboardAvoidingView,
    Button,
    Alert,
    View,
    Text,
    ActivityIndicator
} from 'react-native';
import {fetch} from 'react-native-ssl-pinning';

const AuthScreen = (props) => {
    const [isLoading, setIsLoading] = useState(false);
    const [email, onChangeEmail] = useState('');
    const [password, onChangePassword] = useState('');
    // redirect to home if already logged in
    // if (authenticationService.currentUserValue) {
    //     props.navigation.navigate({ routeName: 'HomePage' });
    // }
    const handleButtonPressed = async () => {
        setIsLoading(true);
        // authenticationService.login(email, password)
        //     .then(
        //         data => {
        //             props.navigation.navigate({ routeName: 'HomePage' });
        //         },
        //         error => {
        //             setIsLoading(false);
        //             console.log("Something went wrong");
        //         }
        // );
        fetch(
            'https://192.168.1.22:8080/user/login',
            {
                method: 'POST',
                headers: {
                    "Content-type": "application/json; charset=UTF-8"
                },
                body: JSON.stringify({
                    user_id: email,
                    user_password: password,
                }),
                sslPinning: {
                    certs: ["cert1"], // cert file name without the `.cer`
                  }
            }
        ).then(function(res) { 
            setIsLoading(false)
            if(res.status === 200) {
                props.navigation.navigate({ routeName: 'HomePage' });
            } else {
                console.log("Error");
            } 
            //console.log(`We got your response! Response - ${JSON.stringify(res)}`),
            })
        .catch(
            setIsLoading(false),
            err => console.log(`Whoopsy doodle! Error - ${err}`)
        )
    }
    return (
        <KeyboardAvoidingView
            behavior="padding"
            keyboardVerticalOffset={50}
            style={styles.screen}>
            <TextInput style={styles.authContainer}
                onChangeText={text => onChangeEmail(text)}
                value={email}
            />
            <TextInput style={styles.authContainer}
                onChangeText={password => onChangePassword(password)}
                value={password}
            />
            <View style={styles.buttonContainer}>
                {isLoading ? (<ActivityIndicator size="small" />)
                    : (<Button title="Login" onPress={handleButtonPressed} />)}
            </View>
        </KeyboardAvoidingView>
    );
}

const styles = StyleSheet.create({
    screen: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center'
    },
    authContainer: {
        width: '80%',
        maxWidth: 400,
        maxHeight: 400,
        padding: 20,
        borderColor: 'gray',
        borderWidth: 1,
        marginVertical: 8
    },
    buttonContainer: {
        marginTop: 16,
        paddingVertical: 8,
        //borderWidth: 4,
        borderColor: "#20232a",
        borderRadius: 6,
        backgroundColor: "white",
        color: "#20232a",
        textAlign: "center",
        fontSize: 30,
        fontWeight: "bold"
    }
})

export default AuthScreen;