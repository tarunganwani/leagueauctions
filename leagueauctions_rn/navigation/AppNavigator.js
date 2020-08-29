import React from 'react';
import { createAppContainer, createSwitchNavigator } from 'react-navigation';
import { createStackNavigator } from 'react-navigation-stack';
import AuthScreen from '../screens/users/AuthScreen';
import Register from '../screens/users/Register';
import VerifyOtp from '../screens/users/VerifyOtp'
import FirstScreen from '../screens/users/FirstScreen';
import EditProfile from '../screens/users/EditProfile';
import AuctionsNavigator from './AuctionsNavigator';

const AuthNavigator = createStackNavigator({
    StartPage: FirstScreen,
    LoginPage: AuthScreen,
    RegisterPage: Register,
    OtpPage: VerifyOtp,
    EditProfilePage: EditProfile
}, {
    initialRouteName: 'StartPage'
});

const MainNavigator = createSwitchNavigator({
    AuthPage: AuthNavigator,
    AppPage: AuctionsNavigator,
}, {
    initialRouteName: 'AuthPage'
});


export default createAppContainer(MainNavigator);