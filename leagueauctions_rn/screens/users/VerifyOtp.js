import React, { useState, useRef, useEffect } from 'react';
import {
    SafeAreaView,
    View,
    Text,
    StyleSheet,
    TextInput,
    TouchableOpacity,
    Dimensions
} from 'react-native';
import OTPInputView from '@twotalltotems/react-native-otp-input'
const { width: WIDTH } = Dimensions.get('window')

const VerifyOtp = (props) => {

    return (
        <SafeAreaView style={styles.container}>
            <Text>Enter OTP sent to your EMAIL ID</Text>
            <OTPInputView
                style={styles.otpInput}
                pinCount={6}
                placeholderTextColor={'blue'}
                // code={this.state.code} //You can supply this prop or not. The component will be used as a controlled / uncontrolled component respectively.
                // onCodeChanged = {code => { this.setState({code})}}
                autoFocusOnLoad
                codeInputFieldStyle={styles.underlineStyleBase}
                codeInputHighlightStyle={styles.underlineStyleHighLighted}
                onCodeFilled={(code => {
                    console.log(`Code is ${code}, you are good to go!`)
                })}

            />
            <View style={styles.buttonContainer}>
                <TouchableOpacity onPress={() => { props.navigation.navigate('EditProfilePage') }}>
                    <Text style={styles.text}>Submit</Text>
                </TouchableOpacity>
            </View>
        </SafeAreaView>
    );
}
const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: 'white',
        alignItems: 'center',
        justifyContent: 'center',
    },
    otpInput: {
        width: '80%',
        height: 200,
    },
    borderStyleBase: {
        width: 30,
        height: 45
    },

    borderStyleHighLighted: {
        borderColor: "#03DAC6",
    },

    underlineStyleBase: {
        width: 30,
        height: 45,
        borderWidth: 0,
        borderBottomWidth: 1,
    },

    underlineStyleHighLighted: {
        borderColor: "#03DAC6",
    },
    buttonContainer: {
        width: WIDTH - 300,
        height: 45,
        borderRadius: 25,
        backgroundColor: "#24a0ed",
        justifyContent: 'center',
        marginTop: 20,
    },
    text: {
        fontSize: 16,
        textAlign: 'center'
    },
});

export default VerifyOtp;
