import React, { useState } from 'react';
import {
    TextInput,
    StyleSheet,
    KeyboardAvoidingView,
    Button,
    Alert,
    View,
    Text,
    ActivityIndicator,
    Dimensions,
    TouchableOpacity
} from 'react-native';
import * as authActions from '../../store/actions/auth';
import { useDispatch } from 'react-redux';
const { width: WIDTH } = Dimensions.get('window')
import Icon from 'react-native-vector-icons/Ionicons';
//import Icon from '@expo/vector-icons';

const AuthScreen = (props) => {
    const [isLoading, setIsLoading] = useState(false);
    const [email, onChangeEmail] = useState('');
    const [password, onChangePassword] = useState('');
    const [error, setError] = useState();
    const dispatch = useDispatch();

    const handleButtonPressed = async () => {
        let action = authActions.login(email, password);
        setError(null);
        setIsLoading(true);
        try {
            await dispatch(action);
            props.navigation.navigate('AppPage');
        } catch (err) {
            setError(err.message);
            setIsLoading(false);
        }
    };

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
            <View style={styles.buttonContainer}>
                {isLoading ? (<ActivityIndicator size="small" />)
                    : (
                        <TouchableOpacity onPress={handleButtonPressed} >
                            <Text style={styles.text}>Login</Text>
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
    text:{
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

export default AuthScreen;