import React from 'react';
import Body from './systemEngine'

export default class System extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      planets: [],
      x: this.props.x || 50,
      y: this.props.y || 50,
    }
    this.addPlanet(this.props.planets);
  }

  addPlanet(n) {
    n = n || 1;
    var newPlanets = this.state.planets;
    for (var i = 0; i < n; i++) {
      newPlanets.push(<Body 
        />);
    }
    this.setState(() => ({planets: newPlanets}));
  }

  render() {
    var style = {
      position: 'absolute',
      bottom: this.state.x+'vh',
      right: this.state.y+'vw',
    }

    return(
      <div className='system' style={style}>
        <Body 
          radius={100}
          tilt={0}
          spin={Math.PI / 2000}
          numberOfArcs={10}
          arcTiltIncrement={Math.PI / 10}
          arcBorderWidth={1}
          arcBorderColor={'yellow'}
          arcFillColor={'yellow'}
          arcFillOpacity={5}
          eccentricity={0}
          semimajorAxis={0}
          inclination={0}
          longitudeOfAscendingNode={0}
          argumentOfPeriapsis={0}
          trueAnomaly={0}
        />
        {this.state.planets.map(planet => planet)}
      </div>
    );
  }
}
