import React from 'react';
import { createAppContainer } from 'react-navigation';
import { createStackNavigator } from 'react-navigation-stack';
import AuthScreen from '../components/AuthScreen';
import TestScreen from '../components/TestScreen';

const LoginNavigation = createStackNavigator({
    LoginPage: AuthScreen,
    HomePage: TestScreen
});

export default createAppContainer(LoginNavigation);