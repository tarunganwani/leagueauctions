import React from 'react';
import { createAppContainer, createSwitchNavigator } from 'react-navigation';
import { createStackNavigator } from 'react-navigation-stack';
import AuthScreen from '../screens/users/AuthScreen';
import Register from '../screens/users/Register';
import VerifyOtp from '../screens/users/VerifyOtp'
import FirstScreen from '../screens/users/FirstScreen';
import Profile from '../screens/users/Profile';
import TestScreen from '../screens/users/TestScreen';

const AppNavigator = createStackNavigator({
    HomePage: TestScreen
});

const AuthNavigator = createStackNavigator({
    StartPage: FirstScreen,
    LoginPage: AuthScreen,
    RegisterPage: Register,
    OtpPage: VerifyOtp,
    ProfilePage: Profile
}, {
    initialRouteName: 'StartPage'
});

const MainNavigator = createSwitchNavigator({
    AuthPage: AuthNavigator,
    AppPage: AppNavigator
}, {
    initialRouteName: 'AuthPage'
});


export default createAppContainer(MainNavigator);