function animation() {
  drawSky('System', 'rgba(255,255,255,0.75)', 500, 2, 1, 2);
  spin('star', 0.01);
  spin('orangeCore', 0.02);
  spin('redCore', 0.03);
  spin('green', 0.05);
  spin('blue', 0.05);
  spin('red', 0.05);
  spin('violet', 0.05);
  orbit('star', 0, 0, 0, 0, 0, time);
  orbit('orangeCore', 0, 0, 0, 0, 0, time);
  orbit('redCore', 0, 0, 0, 0, 0, time);
  orbit('green', 150, 0, 0, 0, time);
  orbit('blue', -200, -200, 0, 1.5, time);
  orbit('red', -200, -200, 0, 1.5, time);
  orbit('violet', -450, -150, 0, -1, time);
  time += 1;
}

setUpBody('star', 200, 0, 6, 30, 1, 'solid', '#fffa6b', 'rgba(227,255,68,0.01)');
setUpBody('orangeCore', 150, 0, 12, 30, 1, 'solid', '#e5a159', 'rgba(255, 104, 81, 0.01)');
setUpBody('redCore', 50, 0, 24, 30, 1, 'solid', '#ff6851', 'rgba(227,255,68,0)');
setUpBody('green', 20, 50, 8, 0, 1, 'solid', '#42f453', 'rgba(76, 76, 76 ,0)');
setUpBody('blue', 50, 0, 12, 50, 1, 'solid', 'rgb(191, 201, 255)', 'rgba(150, 166, 255, 0.25)');
setUpBody('red', 100, 0, 7, 50, 1, 'solid', 'rgb(66, 61, 60)', 'rgba(66, 61, 60, 0.01)');
setUpBody('violet', 30, 0, 20, 40, 1, 'solid', 'rgb(255, 211, 243)', 'rgba(255, 211, 243, 0)');

var time = 0;
var systemAnimation = setInterval(animation, 10);
