import React, {Component} from "react";
import Container from '@material-ui/core/Container';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Grid from '@material-ui/core/Grid';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Typography from '@material-ui/core/Typography';
import Accordion from '@material-ui/core/Accordion';
import AccordionSummary from '@material-ui/core/AccordionSummary';
import AccordionDetails from '@material-ui/core/AccordionDetails';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import * as axios from 'axios';

const styles = makeStyles((theme) => ({
    container: {
        marginTop: "20px",
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 240,
    },
    selectEmpty: {
        marginTop: theme.spacing(2),
    },
    heading: {
        fontSize: theme.typography.pxToRem(15),
        fontWeight: theme.typography.fontWeightRegular,
    },
}));

class ComponentsPage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            repositoryContext: this.defaultRepositoryContexts[0],
            components: [
                'my-component',
            ],
            componentToVersion: {
                'my-component': ["0.0.1"],
            }
        }
    }

    componentDidMount() {
        axios({
            method: "post",
            url: "http://localhost:8080/listComponents",
            data: {
                repositoryContext: this.state.repositoryContext,
            }
        }).then((res) => {
            this.setState({
                ...this.state,
                components: res.data.components,
            })
            console.log(res);
        })
    }

    handleChange = (event, newValue) => {
        this.setState({
            repositoryContext: newValue,
            ...this.state
        })
    }

    handleComponentDetailsExpand = (componentName) => (event, expanded) => {
        if (!expanded) {
            return;
        }
        axios({
            method: "post",
            url: "http://localhost:8080/listComponentVersions",
            data: {
                repositoryContext: this.state.repositoryContext,
                componentName: componentName
            }
        }).then((res) => {
            console.log(res)
            const state = {
                ...this.state,
            }
            state.componentToVersion[componentName] = res.data.versions
            this.setState(state)
        })
    }

    getComponentDetails(componentName) {
        if (!this.state.componentToVersion[componentName]) {
            return "No data";
        }
        return (
            <div>
                {this.state.componentToVersion[componentName].map(version => {
                    return <Typography variant="body1">
                        {version}
                    </Typography>
                })}
            </div>
        )
    }

    defaultRepositoryContexts = [
        "eu.gcr.io/gardener-project/development",
    ]

    render() {
        const { classes } = this.props;
        return (
            <Container className={classes.container}>
                <Grid container>
                    <Grid item xs={6}>
                        <Breadcrumbs aria-label="breadcrumb">
                            {this.state.repositoryContext.split("/").map(value => {
                                return <Typography color="textPrimary">{value}</Typography>
                            })}
                        </Breadcrumbs>
                    </Grid>
                    <Grid item>
                        <Autocomplete
                            id="combo-box-demo"
                            options={this.defaultRepositoryContexts}
                            getOptionLabel={(option) => option}
                            onChange={this.handleChange}
                            style={{width: 300}}
                            renderInput={(params) => <TextField {...params} label="Repository Context"
                                                                variant="outlined"/>}
                        />
                    </Grid>
                </Grid>

                {this.state.components.map((row) => {
                    return <Accordion onChange={this.handleComponentDetailsExpand(row)}>
                        <AccordionSummary
                            expandIcon={<ExpandMoreIcon />}
                            aria-controls="comp-header-{row}"
                            id="comp-header-{row}"
                            >
                            <Typography className={classes.heading}>{row}</Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                            {this.getComponentDetails(row)}
                        </AccordionDetails>
                    </Accordion>
                })}
            </Container>
        )
    }
}

export default withStyles(styles, { withTheme: true })(ComponentsPage);