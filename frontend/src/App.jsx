import React from 'react';
import AppBar from '@material-ui/core/AppBar'
import Button from '@material-ui/core/Button';
import PropTypes from 'prop-types';
import {IconButton, Toolbar, Typography, Menu, MenuItem, Tabs, Tab, Box} from "@material-ui/core";
import {Menu as MenuIcon} from "@material-ui/icons";
import HomeIcon from '@material-ui/icons/Home';
import DescriptionIcon from '@material-ui/icons/Description';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import { BrowserRouter as Router, NavLink, Link as RouterLink, Switch, Route, Link, useHistory } from 'react-router-dom';

import NotFoundPage from "./pages/404";
import ComponentsPage from "./pages/components";
import OverviewerPage from "./pages/overviewer";
import HomePage from "./pages/home";

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    menuButton: {
        marginRight: theme.spacing(2),
    },
    link: {
        padding: "6px 16px",
        color: theme.palette.text.primary,
    },
    filler: {
        flexGrow: 1,
    },
    list: {
        width: 250,
    },
    fullList: {
        width: 'auto',
    },
}));

function App() {
    const classes = useStyles();
    const history = useHistory();
    const [state, setState] = React.useState({
        navbar: false,
    });

    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }

        setState({ ...state, navbar: open });
    };

    const doNav = (path) => () => {
        history.push(path);
        setState({ ...state, navbar: false });
    }

    return (
      <React.Fragment>
          <AppBar position="static">
              <Toolbar>
                  <IconButton edge="start" className={classes.menuButton} color="inherit" aria-haspopup="true" onClick={toggleDrawer(true)} aria-label="menu">
                      <MenuIcon />
                  </IconButton>
                  <Typography variant="h6">
                      Landscaper
                  </Typography>
                  <div className={classes.filler} > </div>
                  <Button color="inherit">Login</Button>
              </Toolbar>
          </AppBar>
          <Drawer open={state["navbar"]} onClose={toggleDrawer(false)}>
              <List className={classes.list}>
                  <ListItem button key="Home" onClick={doNav("/")}>
                      <ListItemIcon><HomeIcon /></ListItemIcon>
                      <ListItemText primary="Home"></ListItemText>
                  </ListItem>
                  <ListItem button key="Components" onClick={doNav("/components")}>
                      <ListItemIcon><DescriptionIcon /></ListItemIcon>
                      <ListItemText primary="Components"></ListItemText>
                  </ListItem>
                  <ListItem button key="Overviewer" onClick={doNav("/overviewer")}>
                      <ListItemIcon><DescriptionIcon /></ListItemIcon>
                      <ListItemText primary="Overviewer"></ListItemText>
                  </ListItem>
              </List>
          </Drawer>
          <Switch>
                <Route exact path="/components" component={ComponentsPage} />
                <Route exact path="/overviewer" component={OverviewerPage} />
                <Route exact path="/" component={HomePage} />
                <Route exact component={NotFoundPage} />
          </Switch>
      </React.Fragment>
    );
}

export default App;
