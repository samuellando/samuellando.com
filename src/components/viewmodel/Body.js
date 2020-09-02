import BodyModel from '@/components/model/Body.js'

export default {
    name: 'Body',

    props: {
        radius: Number,
        tilt: Number,
        spin: Number,
        numberOfArcs: Number,
        arcTiltIncrement: Number,
        arcBorderWidth: String,
        arcBorderStyle: String,
        arcBorderColor: String,
        arcFillColor: String,
        arcFillOpacity: Number,

        eccentricity: Number,
        semimajorAxis: Number,
        inclination: Number,
        longitudeOfAscendingNode: Number,
        argumentOfPeriapsis: Number,
        trueAnomaly: Number,
        mobile: Boolean,
    },

    data() {
        return {
            props: {
                radius: this.radius,
                tilt: this.tilt,
                spin: this.spin,
                numberOfArcs: this.numberOfArcs,
                arcTiltIncrement: this.arcTiltIncrement,
                arcBorderWidth: this.arcBorderWidth,
                arcBorderStyle: this.arcBorderStyle,
                arcBorderColor: this.arcBorderColor,
                arcFillColor: this.arcFillColor,
                arcFillOpacity: this.arcFillOpacity,

                eccentricity: this.eccentricity,
                semimajorAxis: this.semimajorAxis,
                inclination: this.inclination,
                longitudeOfAscendingNode: this.longitudeOfAscendingNode,
                argumentOfPeriapsis: this.argumentOfPeriapsis,
                trueAnomaly: this.trueAnomaly,
                mobile: this.mobile
            },
            TICK_RATE: 1,
            ticker: null,
            style: {},
            arcStyles: []
        };
    },

    mounted() {
        this.model = BodyModel.factory(this.props);
        this.startTick();
    },

   

    methods: {
        isMobile() {
            return (typeof window.orientation !== "undefined") || (navigator.userAgent.indexOf('IEMobile') !== -1);
        },
        startTick() {
            console.log(this.props.mobile);
            if (this.props.mobile || !this.isMobile()) {
                this.ticker = setInterval(this.tick, this.TICK_RATE);
            }
        },
        tick() {
            this.model.tick();
            this.arcStyles = this.getArcStyles();
            this.style = {
                width: this.props.radius,
                height: this.props.radius,
                margin: '0px auto',
                top: (this.model.position[1]-this.props.radius/2)+'px',
                right: (this.model.position[0]-this.props.radius/2)+'px',
                zIndex: Math.floor(this.model.position[2]),
                position: 'absolute',
            };
        },
        getArcStyles() {
            var t = this.model.t;
            var styles = [];
            var angleIncrement = Math.PI / this.model.numberOfArcs;
            var tilt = this.model.tilt;
            var spin = 100 * this.model.spin;
            var radius = this.model.radius + this.model.position[2];
            var tiltIncrement = this.model.arcTiltIncrement;
            var borderColor = this.model.arcBorderColor;
            var borderWidth = this.model.arcBorderWidth;
            var borderStyle = this.model.arcBorderStyle;
            var fillColor = this.model.arcFillColor;
            var opacity = this.model.arcFillOpacity;

            for (var i = 0; i < this.model.numberOfArcs; i++) {
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
                styles.push(style);
            }

            return styles;
        }
    }
}