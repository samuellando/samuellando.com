function drawSky(name, starColor, starCount, drawNow, starMinSize, starMaxSize) {
  var sky = document.getElementById(name);
  if (sky.starCount == null) {
    sky.starCount = drawNow;
  } else {
    sky.starCount += drawNow;
  }
  if (sky.starCount > starCount) {return;}
  for (var i = 0; i < drawNow; i++) {
    var size = (Math.random()*(starMaxSize-starMinSize)+starMinSize);
    var star = document.createElement('DIV');
    star.style.position = 'absolute';
    star.style.zIndex = '-100';
    star.style.top = Math.random()*window.innerHeight-10;
    star.style.left = Math.random()*window.innerWidth-10;
    star.style.width = size;
    star.style.height = size;
    star.style.backgroundColor = starColor;
    sky.appendChild(star);
  }
}

function setUpBody(name, radius, tilt, numberOfArcs, tiltIncrement, borderWidth, borderFill, color, fillColor) {

  // Setup container for wireframe.
  var body = document.getElementById(name);
  /*
   * In order to allow for plannets to pass infornt of eachother in specific
   * orders we need an absolute positioned element, for convinience we can use
   * a container.
   */
  var arcContainer = document.createElement('DIV');
  body.appendChild(arcContainer);
  body.radius = radius; // Stored so it is accesable to other functions.
  body.style.width = radius;
  arcContainer.style.width = radius; // prevents shakeing.
  body.style.margin = '0 auto';
  arcContainer.style.position = 'absolute';
  arcContainer.style.zIndex = '1'; // This will be changed by the orbit function.

  /*
   * In order to to display a sphere, we consider its two dimentional projection
   * each wire in the wireframe being the inrercept of the sphere and a plane
   * its progection is thus an ellipse, considering this as a vector curve, we
   * can eliminate ti z component to obtain the projection of the z=0 plane.
   */
  var angleIncrement = Math.PI/(numberOfArcs-1); // The arcs are all equally spaced.
  for (var i = 0; i < numberOfArcs; i++) {
    // Each arc is its own div.
    var arc = document.createElement('DIV');
    arc.style.margin = '0 auto';
    // To make the wireframes stack.
    if (i > 0) {
      arc.style.marginTop = '-'+(radius+2*borderWidth)+'px';
    }
    arc.style.border = borderWidth+'px '+borderFill+' '+color; // Define border
    arc.style.backgroundColor = fillColor;
    arc.style.borderRadius = '50%'; // Turns it into a conic.
    arc.angle = angleIncrement*i; // Stored to it is accesable to other functions.
    /*
     * This will cause a rotation on the y axis, sin is used because it is positive
     * between zero and pi. The heght will be held constant, and the body can be
     * rotated to obtain roations on any other axis in the xy plane.
     */
    arc.style.width = (Math.sin(arc.angle)*(radius))+'px';
    arc.style.height = radius+'px';
    /*
     * Here the wireframe is tilted, the tilt increment is just to obtain cool effects.
     */
    arc.style.transform = 'rotate('+(tilt+tiltIncrement*i)+'deg)';
    arcContainer.appendChild(arc);
  }
}
function orbit(name, horizontalAmplitude, verticalAmplitude, horizontalPhaseShift, verticalPhaseShift, time) {
  var body = document.getElementById(name);
  var arcContainer = body.childNodes[0];
  // If the body is not obtibing anything.
  if (horizontalAmplitude == 0 && verticalAmplitude == 0) {
    arcContainer.style.marginTop = (window.innerHeight/2-body.radius/2)+'px';
    arcContainer.style.marginLeft = '0px';
    arcContainer.zIndex = 0;
    return;
  }
  var maxDist = Math.sqrt(verticalAmplitude*verticalAmplitude+horizontalAmplitude*horizontalAmplitude);
  /*
   * This used a rule similar to newtons law of gravity to adjust the speed of the orbit
   * for this we scale time. closer bodies will orbit faster.
   */
  var fg = 300/(maxDist*maxDist);
  // Assumeing a circular orbit, we can use simple harmonic motion.
  var posX = horizontalAmplitude*Math.cos(fg*time+horizontalPhaseShift);
  var posY = verticalAmplitude*Math.sin(fg*time+verticalPhaseShift);
  /*
   * Here the body is positioned in the center of the containing div, then translated
   * by its x and y positions.
   */
  arcContainer.style.marginTop = (window.innerHeight/2-body.radius/2-posY)+'px';
  arcContainer.style.marginLeft = posX+'px';
  // Here we check if the z index can be changed, done at maxDist
  var dist = Math.sqrt(posX*posX+posY*posY);
  if (maxDist-dist < 20) {
    if (arcContainer.flipZ) {
      arcContainer.flipZ = false;
      arcContainer.style.zIndex = (-1*parseInt(arcContainer.style.zIndex))+'';
    }
  } else if (maxDist-dist > 20){
    arcContainer.flipZ = true;
  }
}
function spin(name, spinStep) {
  var body = document.getElementById(name);
  var arcs = body.childNodes[0].childNodes;
  for (var i = 0; i < arcs.length; i++) {
    // If the anlge is Pi, reset to zero it better to make a jump that is to small.
    if (arcs[i].angle < Math.PI) {
      arcs[i].angle += spinStep;
    } else {
      arcs[i].angle = 0;
    }
    // Update to width of the elipse.
    arcs[i].style.width = (Math.sin(arcs[i].angle)*(body.radius))+'px';
  }
}

