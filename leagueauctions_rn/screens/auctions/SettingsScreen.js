import React from 'react'
import {
    View,
    Text,
    StyleSheet,
    SafeAreaView,
    TouchableOpacity,
    Dimensions,
    Image
} from 'react-native';
import { Ionicons, MaterialIcons } from "react-native-vector-icons";
const { width: WIDTH } = Dimensions.get('window')

const HomeScreen = () => {
    return (
        <SafeAreaView style={styles.container}>
            <TouchableOpacity>
                <View style={styles.profileInputContainer}>
                    <Image source={require("../../assets/blank.png")} style={styles.image}></Image>
                    <Text style={styles.profileText}> League Auctions </Text>
                    <Text style={styles.profileText}> Right Hand Batsman</Text>
                    <Ionicons name='ios-arrow-forward' size={25} color="#B3B6B7" style={{ position: 'absolute', top: 8, right: 10, paddingTop: 30 }} />
                </View>
            </TouchableOpacity>
            <View style={styles.inputContainer}>
                <TouchableOpacity style={{ borderBottomWidth: 1, borderColor: '#D0D3D4' }}>
                    <View style={styles.boxView}>
                        <Ionicons name='ios-key' size={25} color="#4F8EF7" style={styles.inputIcon} />
                        <Text style={styles.text}> Account</Text>
                        <Ionicons name='ios-arrow-forward' size={25} color="#B3B6B7" style={{ position: 'absolute', top: 8, right: 10 }} />
                    </View>
                </TouchableOpacity>
                <TouchableOpacity>
                    <View style={styles.boxView}>
                        <Ionicons name='ios-information-circle' size={25} color="#4F8EF7" style={styles.inputIcon} />
                        <Text style={styles.text}> Contact Us</Text>
                        <Ionicons name='ios-arrow-forward' size={25} color="#B3B6B7" style={{ position: 'absolute', top: 8, right: 10 }} />
                    </View>
                </TouchableOpacity>
            </View>
            <View style={styles.inputContainer}>
                <TouchableOpacity>
                    <View style={styles.boxView}>
                        <Ionicons name='ios-log-out' size={25} color="#4F8EF7" style={styles.inputIcon} />
                        <Text style={styles.text}> Log Out</Text>
                        <Ionicons name='ios-arrow-forward' size={25} color="#B3B6B7" style={{ position: 'absolute', top: 8, right: 10 }} />
                    </View>
                </TouchableOpacity>
            </View>
        </SafeAreaView>
    );
}
const styles = StyleSheet.create({
    image: {
        //flex: 1,
        height: 80,
        width: 80,
        top: 8,
        left: 10,
        position: "absolute"
    },
    container: {
        flex: 1,
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'stretch',
    },
    inputIcon: {
        position: 'absolute',
        top: 8,
        left: 10
    },
    text: {
        fontSize: 16,
        textAlign: 'left',
        paddingTop: 12,
        paddingLeft: 37,
        color: "#000000",
        fontFamily: 'Palatino-Bold'
    },
    profileText: {
        fontSize: 18,
        textAlign: 'left',
        paddingTop: 15,
        paddingLeft: 100,
        color: "#000000",
        fontFamily: 'Palatino-Bold',
        justifyContent: 'center'
    },
    buttonContainer: {
        width: WIDTH - 300,
        height: 45,
        justifyContent: 'center',
        marginTop: 20,
    },
    boxView: {
        height: 40,
        backgroundColor: 'white',
        // borderColor: 'black',
        // borderBottomWidth: 1,
        // borderTopWidth: 1,
    },
    inputContainer: {
        marginBottom: 50,
        borderColor: '#B3B6B7',
        borderBottomWidth: 1,
        borderTopWidth: 1,
    },
    profileInputContainer: {
        marginBottom: 50,
        borderColor: '#B3B6B7',
        borderBottomWidth: 1,
        borderTopWidth: 1,
        height: 100
    },
});

export default HomeScreen