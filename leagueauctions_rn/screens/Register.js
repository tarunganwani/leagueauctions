import React, { useState,useRef } from 'react'
import { StyleSheet, Text, KeyboardAvoidingView, Button, TextInput, TouchableOpacity, View } from 'react-native';
//simport { useForm } from 'react-hook-form';

const Register = props => {
    const handleButtonPressed = () => {
          // const handleButtonPressed = () => {
    //     const requestOptions = {
    //         "user_id": username,
    //         "user_password": password
    //     };

    //     return axios.post(`https://localhost:8081/user/login`, requestOptions)
    //         .then(function (response) {
    //             console.log(response);
    //             if (response.status === 200) {
    //                 return <OtpScreen />
    //             } else {
    //                 console.log("Some error ocurred " + response.data);
    //             }
    //         })
    //         .catch(function (error) {
    //             console.log(error);
    //         })
        props.navigation.navigate('HomePage');
    }

    
    
    
    const passwordFocus = useRef();
    const first = useRef();
    const confirmPasswordFocus = useRef();

    const [name, setName] = useState('');
    const [password, setPassword] = useState('');

    
    const onButtonPressed = (val) => {
       // console.warn('inside function');
        if(val=='mayur')
            alert('valid');
    
    }
    

    return (
        <KeyboardAvoidingView style={styles.screen}>



            
            <TextInput
                autoCapitalize="none"
                placeholder='Your email'
                style={styles.textInputStyle}
                keyboardType="email-address"
                onChangeText={(val) => setName(val)}
                //autoFocus={true}
                //returnKeyType="next"
                ref={first}
                onSubmitEditing={() => passwordFocus.current.focus()}
                blurOnSubmit={false}

               
            />
            
            <TextInput
                placeholder='Password'
                secureTextEntry={true}
                style={styles.textInputStyle}
                ref={passwordFocus}
                onSubmitEditing={() => confirmPasswordFocus.current.focus()}
                onChangeText={(val) => setPassword(val)}
                onKeyPress={(e) => {e.nativeEvent.key=='Backspace' ? password=="" ? first.current.focus() : undefined : undefined}} 
                blurOnSubmit={false}
                //ref={register({ required: true })}
                
                
            />

            <TextInput
                placeholder='Confirm Password'
                secureTextEntry={true}
                style={styles.textInputStyle}
                ref={confirmPasswordFocus}
            />


            <TouchableOpacity 
            style={styles.button}
            onPress={onButtonPressed({name})}>
                <Text>Sign Up </Text>
            </TouchableOpacity>


            <Text>
                name:{name},
                password:{password}
            </Text>



        </KeyboardAvoidingView>
    );
}



const styles = StyleSheet.create({
    screen: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center'
    },

    /*firstview: {
      flexDirection: 'column',  
      justifyContent: 'center'
    },*/

    textInputStyle: {
        marginBottom: 40,
        width: 200,
        borderColor: 'black',
        borderWidth: 2,
        padding: 10,
        paddingBottom: 5
    },

    button: {
        //alignSelf: 'stretch',
        alignItems: 'center',
        padding: 20,
        backgroundColor: '#808080',
        marginBottom: 20,
        width: 200,
        //fontWeight: 'bold'
       // height: 70,
        //paddingStart: 0
    
      }
});

export default Register;