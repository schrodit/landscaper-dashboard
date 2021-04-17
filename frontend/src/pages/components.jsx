import React, {Component} from "react";
import Container from '@material-ui/core/Container';
import { makeStyles, withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Autocomplete from '@material-ui/lab/Autocomplete';
import Paper from '@material-ui/core/Paper';
import Grid from '@material-ui/core/Grid';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Typography from '@material-ui/core/Typography';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
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
}));

class ComponentsPage extends Component {

    constructor(props) {
        super(props);
        this.state = {
            repositoryContext: this.defaultRepositoryContexts[0],
            components: [
                { name: 'my-component', version: '0.0.1' },
            ]
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
    };

    defaultRepositoryContexts = [
        "eu.gcr.io/gardener-project/development",
    ];

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
                <TableContainer component={Paper}>
                    <Table className={classes.table} aria-label="simple table">
                        <TableHead>
                            <TableRow>
                                <TableCell>Name </TableCell>
                                <TableCell align="right">Version</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {this.state.components.map((row) => (
                                <TableRow key={row.name}>
                                    <TableCell component="th" scope="row">
                                        {row.name}
                                    </TableCell>
                                    <TableCell align="right">{row.version}</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Container>
        )
    }
}

export default withStyles(styles, { withTheme: true })(ComponentsPage);