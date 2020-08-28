export default {
    factory(props) {
        return {
            position: [0, 0, 0],
            t: 0,

            radius: or(props.radius, Math.random() * 100 + 40),
            tilt: or(props.tilt, Math.random() * (Math.PI / 2)),
            spin: or(props.spin, Math.random() * Math.PI / 100),
            numberOfArcs: or(props.numberOfArcs, Math.floor(Math.random() * 20 + 1)),
            arcTiltIncrement: or(props.arcTiltIncrement, Math.random() * Math.PI/2),
            arcBorderWidth: or(props.arcBorderWidth, Math.random()*2),
            arcBorderStyle: or(props.arcBorderStyle, 'solid'),
            arcBorderColor: or(props.arcBorderColor, getRandomColor()),
            arcFillColor: or(props.arcFillColor, getRandomColor()),
            arcFillOpacity: or(props.arcFillOpacity, 5),

            eccentricity: or(props.eccentricity, Math.random()*0.5),
            semimajorAxis: or(props.semimajorAxis, Math.random()*200+20),
            inclination: or(props.inclination, Math.random()*Math.PI/2),
            longitudeOfAscendingNode: or(props.longitudeOfAscendingNode, Math.random()*Math.PI/2),
            argumentOfPeriapsis: or(props.argumentOfPeriapsis, Math.random()*Math.PI/2),
            trueAnomaly: or(props.trueAnomaly, Math.random()*2*Math.PI),

            getPosition() {
                var a = this.semimajorAxis;
                var b = (1.0 - this.eccentricity)*a;
                var t = this.t;
                var t0 = this.trueAnomaly;
                var i = this.inclination;
                var O = this.longitudeOfAscendingNode;
                var w = this.argumentOfPeriapsis;
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
                var matrix = matrixMultiply(rz,rx,ry,p);
                return [matrix[0][0], matrix[1][0], matrix[2][0]];
            },

            tick() {
                var x = this.position[0];
                var y = this.position[1];
                var z = this.position[2];

                // This will relate the speed with the distace from the center.
                if (x + y + z != 0) {
                    this.t += 200/(x*x+y*y+z*z);
                } else {
                    this.t += 0.02;
                }

                this.position = this.getPosition(this.t);
            }
        }
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