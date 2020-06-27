import React, { useState,useRef, Component } from 'react'
import { StyleSheet, Text, KeyboardAvoidingView, Button, TextInput, TouchableOpacity, View } from 'react-native';
import { NavigationContainer } from 'react-navigation'


class RegisterDummy extends Component {

    constructor() {

        super() 
        this.state = {
            email:"",
            password:"",
            confirmPassword:""
        }
    }
validate = () => {

                const { email, password, confirmPassword } = this.state
                let reg = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
                if (reg.test(email) === false) 
                {
                    //name : 'Please fill details'
                    //alert("please fill")
                    this.setState({Err: 'Please correct the name'});
                    return false
                }
                else if(password.length < 2)
                {   
                    alert('incorrect');
                    return false;
                }
                else
                    if(password != confirmPassword) {
                        this.setState({ErrforPassword : 'Both the passwords are not same'})
                        return false;
                
                    }
                    else {
                        return true;        
                    }
                        

            }

    render () {
        const handleButtonPressed = () => {

            
            if(this.validate())
            {


            //     const requestOptions = {
            //         "user_id": username,
            //         "user_password": password
            //     };
            
            //     return axios.post(`https://localhost:8081/user/register`, requestOptions)
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
            //navigator.navigate('Homepage');
              //this.navigation.navigate('HomePage');
              //this.navigator.navigate('HomePage');
              this.props.navigation.navigate('HomePage');
              
            }}

            


        return(
            <View style={styles.screen}>
                

                <Text style={styles.title}>Auctions</Text>

                
                <TextInput
                autoCapitalize="none"
                placeholder='Your email'
                style={styles.textInputStyle}
                keyboardType="email-address"
                onChangeText={(val) => this.setState({email : val})}
                //autoFocus={true}
                //returnKeyType="next"
                //onSubmitEditing={() => passwordFocus.current.focus()}
                blurOnSubmit={false}
                
               
            />

            <Text style= {{ color: 'red', textAlign: 'center' }}>
                { this.state.Err }
            </Text>


            <TextInput
                placeholder='Password'
                secureTextEntry={true}
                style={styles.textInputStyle}
                //ref={passwordFocus}
                //onSubmitEditing={() => confirmPasswordFocus.current.focus()}
                onChangeText={(val) => this.setState({password : val})}
                blurOnSubmit={false}
                //ref={register({ required: true })}
                
                
            />

            

            <TextInput
                placeholder='Confirm Password'
                secureTextEntry={true}
                style={styles.textInputStyle}
                onChangeText={(val) => this.setState({confirmPassword : val})}
                //ref={confirmPasswordFocus}
            />

            <Text style= {{ color: 'red', textAlign: 'center' }}>
                { this.state.ErrforPassword }
            </Text>
            
            <TouchableOpacity 
            style={styles.button}
            onPress={handleButtonPressed} >
                <Text>Sign Up </Text>
            </TouchableOpacity>


            <Text>
                name:{this.state.email},
                password:{this.state.password}
            </Text>


            </View>
        );
    }

}



//function handle(){

 
    //const passwordFocus = useRef();
    //const confirmPasswordFocus = useRef();

    //const [name, setName] = useState('');
    //const [password, setPassword] = useState('');

const styles = StyleSheet.create({
    title: {
        marginTop: 16,
        paddingVertical: 8,
        borderWidth: 4,
        borderColor: "#20232a",
        borderRadius: 6,
        backgroundColor: "#61dafb",
        color: "#20232a",
        textAlign: "center",
        fontSize: 30,
        fontWeight: "bold"
      },
    screen: {
        flex: 1,
        //flexDirection: 'column',
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


export default RegisterDummy;