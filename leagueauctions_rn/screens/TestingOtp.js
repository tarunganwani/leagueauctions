import React, { useState,useRef,Component } from 'react'
import { StyleSheet, Text, KeyboardAvoidingView, Button, TextInput, TouchableOpacity, View } from 'react-native';
import OTPInputView from '@twotalltotems/react-native-otp-input'


class TestingOtp extends Component {
    constructor() {

        super() 
        this.state = {
            email:"",
            ButtonStateHolder : true
        }
    }
    enableButton = () => {
        this.setState ({
            //semail:"",
            ButtonStateHolder : false
        })

    }
    validate = () => {
        const { email } = this.state
        if(email.length != 6 ) {
            this.enableButton();
            //this.ButtonStateHolder = false;
        //this.Button.disabled={state.disabled={false}};
        //this.props.Button.disabled={false};
        //props.Button.disabled={false}
        }
    }
 render(){
    return (
        <View>
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
               // onChangeText={this.validate()}
               
            />
            <TouchableOpacity  
            style={styles.button}
            disabled={this.validate()}
                        >
                            <Text>Hellos</Text>
                        </TouchableOpacity>
            
        </View>
    );
}
}
const styles = StyleSheet.create({
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
export default TestingOtp;