import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import AuthScreen from './components/AuthScreen';
import TestScreen from './components/TestScreen';
import LoginNavigation from './navigation/LoginNavigation';

export default function App() {
  return <LoginNavigation />
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#191',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
