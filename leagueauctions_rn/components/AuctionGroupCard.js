import React from 'react';
import { StyleSheet, Text, View ,Image} from 'react-native';
import { Avatar, Button, Card, Title, Paragraph} from 'react-native-paper';
import {Ionicons} from '@expo/vector-icons';

export default function AuctionGroupCards(props){
          return (
               <View style={{margin:5}}>
                    <Card style={{elevation:5}}>
                        <View style={{flexDirection:"row",padding:10}}>
                            <View>
                               <Image
                                style={{height:80,width:80,borderRadius:40}}
                                source={{uri:'https://source.unsplash.com/RDcEWH5hSDE/600x500'}}
                                 />
                            </View>
                            <View>
                                  <Card.Content>
                                    <Title>{props.name}</Title>
                                    <Paragraph >hello coders never quit</Paragraph>
                                    </Card.Content>
                            </View>
                         
                        </View>
                      
                     </Card>
                   </View>
               
            );
}