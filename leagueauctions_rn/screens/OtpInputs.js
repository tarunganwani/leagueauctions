import React, { useState,useRef, Component } from 'react'
import { StyleSheet, Text, KeyboardAvoidingView, Button, TextInput, TouchableOpacity, View } from 'react-native';
//import OtpInputs from 'react-native-otp-inputs';
import { Content, Item, Input } from 'native-base';
import { Grid, Col } from 'react-native-easy-grid';

class OtpInputs extends Component {
   
  state={otp:[]};
  otpTextInput = [];

  componentDidMount() {
      this.otpTextInput[0]._root.focus();
  }

  renderInputs() {
      const inputs = Array(6).fill(0);
      const txt = inputs.map(
          (i, j) => <Col key={j} style={styles.txtMargin}><Item regular>
              <Input
                  style={[styles.inputRadius, { borderRadius: 10 }]}
                  keyboardType="numeric"
                  onChangeText={v => this.focusNext(j, v)}
                  onKeyPress={e => this.focusPrevious(e.nativeEvent.key, j)}
                  ref={ref => this.otpTextInput[j] = ref}
              />
          </Item></Col>
      );
      return txt;
  }

  focusPrevious(key, index) {
      if (key === 'Backspace' && index !== 0)
          this.otpTextInput[index - 1]._root.focus();
  }

  focusNext(index, value) {
      if (index < this.otpTextInput.length - 1 && value) {
          this.otpTextInput[index + 1]._root.focus();
      }
      if (index === this.otpTextInput.length - 1) {
          this.otpTextInput[index]._root.blur();
      }
      const otp = this.state.otp;
      otp[index] = value;
      this.setState({ otp });
     // this.props.getOtp(otp.join(''));
  }


  render() {
      return (
          <Content padder>
              <Grid style={styles.gridPad}>
                  {this.renderInputs()}
              </Grid>
          </Content>
    );
            }
}

const styles = StyleSheet.create({
  gridPad: { padding: 30 },
  txtMargin: { margin: 3 },
  inputRadius: { textAlign: 'center' }
});

export default OtpInputs;