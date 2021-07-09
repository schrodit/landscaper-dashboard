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
import { CheckCircleOutline, ErrorOutline, FlashOn } from '@material-ui/icons';

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

class OverviewerPage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            installations: [
                {
                    Name: 'inst-0',
                    Namespace: 'dummy-ns',
                    UpToDate: "true",
                    Phase: 'Failed',
                    Subinstallations: [
                        {
                            Name: 'subinst-0-0',
                            Namespace: 'dummy-ns',
                            UpToDate: "true",
                            Phase: 'Successful',
                            Subinstallations: [
                                {
                                    Name: 'subinst-1-0',
                                    Namespace: 'dummy-ns',
                                    UpToDate: "true",
                                    Phase: 'Successful',
                                    Subinstallations: [],
                                },
                            ],
                        },
                        {
                            Name: 'subinst-0-1',
                            Namespace: 'dummy-ns',
                            UpToDate: "true",
                            Phase: 'Failed',
                            Subinstallations: [],
                        }
                    ],
                    Execution: null,
                },
                {
                    Name: 'inst-1',
                    Namespace: 'dummy-ns',
                    UpToDate: "true",
                    Phase: 'Successful',
                    Subinstallations: [],
                    Execution: null,
                },
                {
                    Name: 'inst-2',
                    Namespace: 'dummy-ns',
                    UpToDate: "true",
                    Phase: 'Progressing',
                    Subinstallations: [],
                    Execution: null,
                }
            ],
        }
    }

    componentDidMount() {
        axios({
            method: "post",
            url: "http://localhost:8080/listComponents",
            data: {}
        }).then((res) => {
            this.setState({
                installations: res,
            })
        })
    }

    // handleChange = (event, newValue) => {
    //     this.setState({
    //         repositoryContext: newValue,
    //         ...this.state
    //     })
    // }

    // handleComponentDetailsExpand = (componentName) => (event, expanded) => {
    //     if (!expanded) {
    //         return;
    //     }
    //     axios({
    //         method: "post",
    //         url: "http://localhost:8080/listInstallationData",
    //         data: {}
    //     }).then((res) => {
    //         console.log(res)
    //         const state = {
    //             ...this.state,
    //         }
    //         state.componentToVersion[componentName] = res.data.versions
    //         this.setState(state)
    //     })
    // }

    getPhaseIcon(phase) {
        if (phase === "Successful") {
            return (<CheckCircleOutline color="primary" />)
        } else if (phase === "Failed") {
            return (<ErrorOutline />)
        } else {
            return (<FlashOn />)
        }
    }

    renderInstallationEntryHeader(inst) {
        return <div>{inst.Namespace} / {inst.Name} {this.getPhaseIcon(inst.Phase)}</div>
    }

    renderInstallationDetails(inst) {
        return (
            <div>
                Name: {inst.Name}<br />
                Namespace: {inst.Namespace}<br />
                UpToDate: {inst.UpToDate}<br />
                Phase: {inst.Phase}<br />
                Subinstallations:<br />
                <div>
                    {inst.Subinstallations.map((sinst) => {return this.renderInstallationEntry(sinst)})}
                </div>
                Execution: null
            </div>
        )
    }

    renderInstallationEntry(inst) {
         return (
            <Accordion>
                <AccordionSummary expandIcon={<ExpandMoreIcon />} aria-controls="panel1a-content" id="panel1a-header">
                    <Typography>{this.renderInstallationEntryHeader(inst)}</Typography>
                </AccordionSummary>
                <AccordionDetails>
                    <Typography>{this.renderInstallationDetails(inst)}</Typography>
                </AccordionDetails>
            </Accordion>
        )
    }

    render() {
        const { classes } = this.props;
        return (
            <Container className={classes.container}>
                {this.state.installations.map((inst) => {return this.renderInstallationEntry(inst)})}
            </Container>
        )
    }
}

export default withStyles(styles, { withTheme: true })(OverviewerPage);
