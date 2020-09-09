import React from 'react';
import { createAppContainer, createSwitchNavigator } from 'react-navigation';
import { createStackNavigator } from 'react-navigation-stack';
import HomeScreen from '../screens/auctions/HomeScreen';
import SettingsScreen from '../screens/auctions/SettingsScreen';
import { createBottomTabNavigator } from 'react-navigation-tabs';
import { Ionicons, MaterialIcons } from "react-native-vector-icons";
import Colors from '../constants/Colors';
import { Platform } from 'react-native';

const AuctionScreensNavigator = createStackNavigator({
    Home: HomeScreen
});

const AuctionsNavigator = createBottomTabNavigator(
    {
        HomePage: {
            screen: AuctionScreensNavigator,
            navigationOptions: {
                tabBarLabel: 'Auctions',
                tabBarIcon: tabInfo => {
                    return (
                        <Ionicons name="ios-home" size={25} color={tabInfo.tintColor} />
                    );
                }
            }
        },
        Settings: {
            screen: SettingsScreen,
            navigationOptions: {
                tabBarLabel: 'Settings',
                tabBarIcon: tabInfo => {
                    return (
                        <Ionicons name="ios-settings" size={25} color={tabInfo.tintColor} />
                    );
                }
            }
        }
    },
    {
        tabBarOptions: {
            activeTintColor: Platform.OS === 'android' ? Colors.accent : Colors.primary
        }
    }
);

export default createAppContainer(AuctionsNavigator);