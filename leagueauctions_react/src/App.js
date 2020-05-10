import React, { Component } from 'react';
import Header from './components/Header/Header';
import Loginscreen from './components/LoginScreen/Loginscreen';
import './App.css';
import { BrowserRouter } from 'react-router-dom';

class App extends Component {

  render() {
    return (
      <BrowserRouter>
        <div className="App">
          <Header />
          <Loginscreen style={style} />
        </div>
      </BrowserRouter>
    );
  }
}

const style = {
  margin: 15,
};

export default App;
