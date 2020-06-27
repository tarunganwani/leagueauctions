import React,  {Component} from 'react';
import { StyleSheet, Text, View, Button, TextInput, TouchableOpacity } from 'react-native';


export default class Signup extends React.Component {


    render() {
    return (

        <View style={styles.screen}>
          <View>
            <TextInput
              placeholder= 'Your email' 
              style={ styles.textInputStyle}
              keyboardType="email-address"
             
              >
            </TextInput>
    
    
    
            <TextInput>
              placeholder= 'Password' 
              secureTextEntry={true}
              style={ styles.textInputStyle}
              
              >
            </TextInput>
    
            <TextInput
              placeholder= 'Confirm Password' 
              secureTextEntry={true}
              style={ styles.textInputStyle}
              
              >
    
            </TextInput>
    
            <TouchableOpacity style={styles.button}>
              <Text>Sign Up </Text>  
            </TouchableOpacity>  
    
            
          </View>
        </View>
      );
    }
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
    
      buttonstyle: {
        width: 200,
        height: 100,
        paddingStart: 40
      },
    
      button: {
        alignSelf: 'stretch',
        alignItems: 'center',
        padding: 20,
        backgroundColor: '#808080',
        marginBottom: 20
    
      }
    
    });