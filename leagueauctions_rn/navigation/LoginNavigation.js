import React from 'react';
import { createAppContainer } from 'react-navigation';
import { createStackNavigator } from 'react-navigation-stack';
import AuthScreen from '../components/AuthScreen';
import TestScreen from '../components/TestScreen';
import RegisterDummy from '../components/RegisterDummy';
import Test from '../components/Test';
import TestingOtp from '../components/TestingOtp';

const LoginNavigation = createStackNavigator({
    LoginPage: TestingOtp,
    HomePage: TestScreen,
    RegisterPage: RegisterDummy,
    OtpPage: Test
});

export default createAppContainer(LoginNavigation);