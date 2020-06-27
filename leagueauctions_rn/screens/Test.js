import * as React from 'react';
import { Text, View, StyleSheet, TextInput, TouchableOpacity } from 'react-native';

//import OtpInputs from '../components/OtpInputs';


class Test extends React.Component {


  constructor(props) {
    super(props);
    this.state = {
      firstInput: "",
      secondInput: "",
      thirdInput: ""
    }

    this.focusNextField = this.focusNextField.bind(this);
    this.inputs = {};
  }

  focusNextField(id) {
    this.inputs[id].focus();
  }
  focusPrevious(key, id) {

    if (key == 'Backspace') {
      this.inputs[id].focus();
    }

  }


  render() {

    return (

      <View style={styles.screen}>



        <TextInput
          maxLength={1}
          style={styles.textInputStyle}
          onChange={(val) => {
            this.focusNextField('two');
            this.setState({ firstInput: val })
          }}
          //returnKeyType={ "next" }
          ref={input => {
            this.inputs['one'] = input;
          }}
        >




        </TextInput>
        <TextInput
          style={styles.textInputStyle}
          onChangeText={(val) => {
            this.focusNextField('three');
            this.setState({ secondInput: val })
          }}
          //returnKeyType={ "next" }
          ref={input => {
            this.inputs['two'] = input;
          }}

          onKeyPress={(e) => { e.nativeEvent.key == 'Backspace' ? this.state.secondInput == "" ? this.focusNextField('one') : undefined : undefined }}
        >




        </TextInput>
        <TextInput
          style={styles.textInputStyle}
          onChangeText={() => {
            this.focusNextField('four');
          }}
          returnKeyType={"next"}
          ref={input => {
            this.inputs['three'] = input;
          }}

        >




        </TextInput>
        <TextInput
          style={styles.textInputStyle}
          onChangeText={() => {
            this.focusNextField('five');
          }}
          returnKeyType={"next"}
          ref={input => {
            this.inputs['four'] = input;
          }}
        >




        </TextInput>
        <TextInput
          style={styles.textInputStyle}
          onChangeText={() => {
            this.focusNextField('six');
          }}
          returnKeyType={"next"}
          ref={input => {
            this.inputs['five'] = input;
          }}
        >




        </TextInput>
        <TextInput
          style={styles.textInputStyle}
          returnKeyType={"done"}

          ref={input => {
            this.inputs['six'] = input;
          }}
        >




        </TextInput>
        <TouchableOpacity
          style={styles.button}
        >
          <Text>Sign Up </Text>
        </TouchableOpacity>

      </View>
    );
  }

}

const styles = StyleSheet.create({
  screen: {
    flex: 0.6,
    flexDirection: 'row',
    justifyContent: 'space-evenly',
  },
  textInputStyle: {
    //flex: 1,
    //justifyContent: 'center',
    //alignItems: 'center',
    //flex: 0.6,
    //flexDirection: 'row',
    //justifyContent: 'space-evenly',
    alignSelf: 'center',
    //flexDirection:'column',
    marginBottom: 40,
    marginLeft: 260,

    width: 50,
    borderColor: 'black',
    borderWidth: 2,
    padding: 10,
    paddingBottom: 5
  },
  button: {
    //alignSelf: 'stretch',
    //flexDirection: 'column',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#808080',
    marginTop: 400,
    marginRight: 360,
    marginBottom: 50,
    //marginLeft: 5,
    width: 150,
    //fontWeight: 'bold'
    height: 50,
    //paddingStart: 50

  },
  startText: {
    padding: 100
  }

});

export default Test;