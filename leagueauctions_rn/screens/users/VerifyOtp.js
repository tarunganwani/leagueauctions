import React from 'react'
import { View, Text, StyleSheet } from 'react-native';

const VerifyOtp = () => {
    return (
        <View style={styles.container}>
            <Text>Hello World</Text>
        </View>
    );
}
const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#191',
        alignItems: 'center',
        justifyContent: 'center',
    },
});

export default VerifyOtp;