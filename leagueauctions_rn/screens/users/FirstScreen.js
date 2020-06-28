import React, { useState } from 'react';
import {
    Button,
    View,
    StyleSheet,
    Dimensions,
    Text
} from 'react-native';
import { TouchableOpacity } from 'react-native-gesture-handler';
const { width: WIDTH } = Dimensions.get('window')

const FirstScreen = (props) => {
    return (
        <View style={styles.screen}>
            <TouchableOpacity style={styles.buttonContainer} onPress={() => { props.navigation.navigate('LoginPage') }} >
                <Text style={styles.text}>Login</Text>
            </TouchableOpacity>

            <TouchableOpacity style={styles.buttonContainer} onPress={() => { props.navigation.navigate('RegisterPage') }} >
                <Text style={styles.text}>Sign Up</Text>
            </TouchableOpacity>
        </View>
    );
}

const styles = StyleSheet.create({
    screen: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center'
    },
    text:{
        fontSize: 16,
        textAlign: 'center'
    },
    buttonContainer: {
        width: WIDTH - 100,
        height: 45,
        borderRadius: 25,
        backgroundColor: "#24a0ed",
        justifyContent: 'center',
        marginTop: 20
    },
})

export default FirstScreen;