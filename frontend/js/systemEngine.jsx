class Body extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      x: 0,
      y: 0,
      z: 0,

      radius: or(this.props.radius, Math.random() * 100 + 40),
      tilt: or(this.props.tilt, Math.random() * (Math.PI / 2)),
      spin: or(this.props.spin, Math.random() * Math.PI / 100),
      numberOfArcs: or(this.props.numberOfArcs, Math.random() * 20 + 1),
      arcTiltIncrement: or(this.props.arcTiltIncrement, Math.random() * Math.PI/2),
      arcBorderWidth: or(this.props.arcBorderWidth, Math.random()*2),
      arcBorderStyle: or(this.props.arcBorderStyle, 'solid'),
      arcBorderColor: or(this.props.arcBorderColor, getRandomColor()),
      arcFillColor: or(this.props.arcFillColor, getRandomColor()),
      arcFillOpacity: or(this.props.arcFillOpacity, 5),

      eccentricity: or(this.props.eccentricity, Math.random()*0.5),
      semimajorAxis: or(this.props.semimajorAxis, Math.random()*200+20),
      inclination: or(this.props.inclination, Math.random()*Math.PI/2),
      longitudeOfAscendingNode: or(this.props.longitudeOfAscendingNode, Math.random()*Math.PI/2),
      argumentOfPeriapsis: or(this.props.argumentOfPeriapsis, Math.random()*Math.PI/2),
      trueAnomaly: or(this.props.trueAnomaly, Math.random()*2*Math.PI),
      t: 0,
    };

    var p = this.getPosition(0);
    this.state.x = p[0][0];
    this.state.y = p[1][0];
    this.state.z = p[2][0];
  }

  tick() {
    var x = this.state.x;
    var y = this.state.y;
    var z = this.state.z;
    var t = this.state.t;

    // This will relate the speed with the distace from the center.
    if (x + y + z != 0) {
      t += 200/(x*x+y*y+z*z);
    } else {
      t += 0.02;
    }

    var p = this.getPosition(t);

    var x = p[0][0];
    var y = p[1][0];
    var z = p[2][0];

    this.setState(() => ({x: x, y: y, z: z, t: t}));
  }

  componentDidMount() {
    this.interval = setInterval(() => this.tick(), 10);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  getPosition(t) {
    var a = this.state.semimajorAxis;
    var b = (1.0 - this.state.eccentricity)*a;
    var t0 = this.state.trueAnomaly;
    var i = this.state.inclination;
    var O = this.state.longitudeOfAscendingNode;
    var w = this.state.argumentOfPeriapsis;
    // Determine the position on the elipse with its focus at origin.
    var p  = matrixAdd([
      [a*Math.cos(t+t0)], 
      [b*Math.sin(t+t0)], 
      [0]], [[a-b], [0], [0]]);
    // Determine the rotations.
    var rz = [
      [Math.cos(O),-1*Math.sin(O),0],
      [Math.sin(O),Math.cos(O),0],
      [0,0,1]
    ];
    var ry = [
      [Math.cos(i),0,-1*Math.sin(i)],
      [0,1,0],
      [Math.sin(i),0,Math.cos(i)]
    ];
    var rx = [
      [1,0,0],
      [0,Math.cos(w),-1*Math.sin(w)],
      [0,Math.sin(w),Math.cos(w)]
    ];
    // Calculate the position.
    return matrixMultiply(rz,rx,ry,p);
  }

  createArcs(t) {
    var arcs = [];
    var angleIncrement = Math.PI / this.state.numberOfArcs;
    var tilt = this.state.tilt;
    var spin = 100 * this.state.spin;
    var radius = this.state.radius;
    var tiltIncrement = this.state.arcTiltIncrement;
    var borderColor = this.state.arcBorderColor;
    var borderWidth = this.state.arcBorderWidth;
    var borderStyle = this.state.arcBorderStyle;
    var fillColor = this.state.arcFillColor;
    var opacity = this.state.arcFillOpacity;

    for (var i = 0; i < this.state.numberOfArcs; i++) {
      var height = radius;
      var width = (Math.abs(Math.sin(angleIncrement*i+spin*t))*(radius));
      var style = {
        height: height+'px',
        width: width+'px',
        position: 'absolute',
        top: (radius - height) / 2+'px',
        right: (radius - width) / 2+'px',
        transform: 'rotate('+(tilt+tiltIncrement*i)+'rad)',
        background: fillColor,
        opacity: opacity+'%',
        borderStyle: borderStyle,
        borderWidth: borderWidth,
        borderRadius: '100%',
        borderColor: borderColor,
      };
      arcs.push(<div className='arc' style={style}></div>);
    }
    return arcs;
  }

  render() {
    var style = {
      width: this.state.radius,
      height: this.state.radius,
      margin: '0px auto',
      bottom: this.state.y+'px',
      left: this.state.x+'px',
      zIndex: this.state.z+'px',
      position: 'absolute',
    };

    var arcs = this.createArcs(this.state.t);

    return(
      <div className='body' style={style}>
        {arcs.map(arc => arc)}
      </div>
    );
  }
}

function matrixAdd() {
  for (var i = 1; i < arguments.length; i++) {
    for (var j = 0; j < arguments[0].length; j++) {
      for (var k = 0; k < arguments[0][0].length; k++) {
        arguments[0][j][k] += arguments[i][j][k];
      }
    }
  }
  return arguments[0];
}

function matrixMultiply() {
  for (var i = 0; i < arguments.length - 1; i++) {
    var j = i+1;
    var m = arguments[i].length;
    var n = arguments[j][0].length;
    var a = Array(m).fill().map(() => Array(n).fill(0)); 
    for (var k = 0; k < m; k++) {
      for (var l = 0; l < n; l++) {
        var dotProduct = 0;
        for (var c = 0; c < m; c++) {
          dotProduct += arguments[i][k][c]*arguments[j][c][l];
        }
        a[k][l] = dotProduct;
      }
    }
    arguments[j] = a;
  }
  return a;
}

function getRandomColor() {
    var letters = '0123456789ABCDEF';
    var color = '#';
    for (var i = 0; i < 6; i++) {
          color += letters[Math.floor(Math.random() * 16)];
        }
    return color;
}

function or(a, b) {
  if (a != undefined) {
    return a;
  } else {
    return b;
  }
}
