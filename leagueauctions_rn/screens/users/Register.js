import React, { useState, useRef, Component } from 'react'
import {
    StyleSheet,
    Text,
    KeyboardAvoidingView,
    Button,
    TextInput,
    TouchableOpacity,
    View,
    Dimensions,
    ActivityIndicator
} from 'react-native';
import * as authActions from '../../store/actions/auth';
import { useDispatch } from 'react-redux';
const { width: WIDTH } = Dimensions.get('window')
import Icon from 'react-native-vector-icons/Ionicons';

// constructor() {

//     super()
//     this.state = {
//         email: "",
//         password: "",
//         confirmPassword: ""
//     }
// }
// validate = () => {

//     const { email, password, confirmPassword } = this.state
//     let reg = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
//     if (reg.test(email) === false) {
//         //name : 'Please fill details'
//         //alert("please fill")
//         this.setState({ Err: 'Please correct the name' });
//         return false
//     }
//     else if (password.length < 2) {
//         alert('incorrect');
//         return false;
//     }
//     else
//         if (password != confirmPassword) {
//             this.setState({ ErrforPassword: 'Both the passwords are not same' })
//             return false;

//         }
//         else {
//             return true;
//         }


// }

// render() {
//     const handleButtonPressed = () => {


//         //if (this.validate()) {


//         //     const requestOptions = {
//         //         "user_id": username,
//         //         "user_password": password
//         //     };

//         //     return axios.post(`https://localhost:8081/user/register`, requestOptions)
//         //         .then(function (response) {
//         //             console.log(response);
//         //             if (response.status === 200) {
//         //                 return <OtpScreen />
//         //             } else {
//         //                 console.log("Some error ocurred " + response.data);
//         //             }
//         //         })
//         //         .catch(function (error) {
//         //             console.log(error);
//         //         })
//         //navigator.navigate('Homepage');
//         //this.navigation.navigate('HomePage');
//         //this.navigator.navigate('HomePage');
//         this.props.navigation.navigate('OtpPage');

//         //}
//     }

const Register = (props) => {
    const [email, onChangeEmail] = useState('');
    const [password, onChangePassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    return (
        <KeyboardAvoidingView
            behavior="padding"
            keyboardVerticalOffset={50}
            style={styles.screen}>
            <View style={styles.inputContainer}>
                <Icon name="ios-person" size={28} color="#4F8EF7" style={styles.inputIcon} />
                <TextInput style={styles.input}
                    placeholder={'Email Id'}
                    placeholderTextColor={'rgba(255, 255, 255, 0.7)'}
                    underlineColorAndroid='transparent'
                    onChangeText={text => onChangeEmail(text)}
                    value={email}
                />
            </View>
            <View style={styles.inputContainer}>
                <Icon name="ios-lock" size={28} color="#4F8EF7" style={styles.inputIcon} />
                <TextInput style={styles.input}
                    placeholder={'Password'}
                    secureTextEntry={true}
                    placeholderTextColor={'rgba(255, 255, 255, 0.7)'}
                    underlineColorAndroid='transparent'
                    onChangeText={password => onChangePassword(password)}
                    value={password}
                />
                <TouchableOpacity style={styles.btnEye}>
                    <Icon name="ios-eye" size={26} color="#4F8EF7" />
                </TouchableOpacity>
            </View>
            <View style={styles.inputContainer}>
                <Icon name="ios-lock" size={28} color="#4F8EF7" style={styles.inputIcon} />
                <TextInput style={styles.input}
                    placeholder={'Confirm Password'}
                    secureTextEntry={true}
                    placeholderTextColor={'rgba(255, 255, 255, 0.7)'}
                    underlineColorAndroid='transparent'
                />
                <TouchableOpacity style={styles.btnEye}>
                    <Icon name="ios-eye" size={26} color="#4F8EF7" />
                </TouchableOpacity>
            </View>
            <View style={styles.buttonContainer}>
                {isLoading ? (<ActivityIndicator size="small" />)
                    : (
                        <TouchableOpacity onPress={() => { props.navigation.navigate('OtpPage') }}>
                            <Text style={styles.text}>Register</Text>
                        </TouchableOpacity>
                    )}
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
    inputContainer: {
        margin: 10,
    },
    text: {
        fontSize: 16,
        textAlign: 'center'
    },
    input: {
        width: WIDTH - 55,
        height: 45,
        borderRadius: 25,
        fontSize: 16,
        paddingLeft: 45,
        backgroundColor: 'rgba(0,0,0,0.35)',
        color: 'rgba(255,255,255,0.7)',
        marginHorizontal: 25
    },
    buttonContainer: {
        width: WIDTH - 300,
        height: 45,
        borderRadius: 25,
        backgroundColor: "#24a0ed",
        justifyContent: 'center',
        marginTop: 20,
        // paddingVertical: 8,
        // //borderWidth: 4,
        // borderColor: "#20232a",
        // borderRadius: 6,
        // backgroundColor: "white",
        // color: "#20232a",
        // textAlign: "center",
        // fontSize: 30,
        // fontWeight: "bold"
    },
    inputIcon: {
        position: 'absolute',
        top: 8,
        left: 37
    },
    btnEye: {
        position: 'absolute',
        top: 8,
        right: 37
    }
})


export default Register;