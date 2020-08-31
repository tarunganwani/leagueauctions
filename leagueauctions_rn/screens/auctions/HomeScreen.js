import React from 'react'
import { View, ScrollView, StyleSheet, Text } from 'react-native';
import AuctionGroupCards from '../../components/AuctionGroupCard';
import { HeaderButtons, Item } from 'react-navigation-header-buttons';
import HeaderButton from '../../components/HeaderButton';

const HomeScreen = () => {
    return (
        <View style={styles.container}>
            <ScrollView>
                <AuctionGroupCards name='Blake' />
                <AuctionGroupCards name='FCL' />
            </ScrollView>
        </View>
    );
}

HomeScreen.navigationOptions = navigationData => {
    return {
        headerRight: () =>
            <HeaderButtons HeaderButtonComponent={HeaderButton}>
                <Item
                    title="Favorite"
                    iconName="ios-star"
                    onPress={() => {
                        console.log('Mark as favorite!');
                    }}
                />
            </HeaderButtons>
    }
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center',
    },
});

export default HomeScreen