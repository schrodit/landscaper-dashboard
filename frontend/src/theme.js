import { createMuiTheme } from '@material-ui/core/styles';
import green from '@material-ui/core/colors/green';

// Create a theme instance.
const theme = createMuiTheme({
    palette: {
        primary: {
            main: '#009688',
        },
        secondary: {
            main: '#19857b',
        },
        error: {
            main: green.A400,
        },
        background: {
            default: '#fff',
        },
    },
});

export default theme;