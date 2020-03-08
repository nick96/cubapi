import React, { FunctionComponent } from 'react';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import {
    Grid, TextField, Button,
    Container, AppBar, Toolbar, Typography,
    CssBaseline, Dialog, DialogContent, DialogTitle, DialogActions, DialogProps,
    Paper, Box
} from '@material-ui/core';
import { createStyles, makeStyles, Theme, createMuiTheme, ThemeProvider } from "@material-ui/core/styles";
import { yellow } from "@material-ui/core/colors";
import { useFormik } from "formik";
import * as Yup from "yup";
import './App.css';
import "typeface-roboto";

const USER_SERVICE_URI = process.env.REACT_APP_USER_SERVICE_URI;

const theme = createMuiTheme({
    palette: {
        primary: yellow,
    },
})

const useStyles = makeStyles((theme: Theme) => createStyles({
    root: {
        flexGrow: 1,
    },
    menuButton: {
        marginRight: theme.spacing(2),
    },
    title: {
        flexGrow: 1,
    },
    backdrop: {
        zIndex: theme.zIndex.drawer + 1,
        color: "#fff",
    },
    submitButton: {
        flexGrow: 1,
    },
    blurb: {
        marginTop: theme.spacing(10),
        margin: theme.spacing(50),
        textAlign: "left"
    }
}));

function App() {
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <div className="App">
                <Router>
                    <Switch>
                        <Route exact path="/">
                            <Home />
                        </Route>
                    </Switch>
                </Router>
            </div>
        </ThemeProvider>
    );
}

interface SignInFormErrors {
    email: string | null,
    password: string | null,
}

interface LoginDialogProps extends DialogProps {
}

const LoginDialog: FunctionComponent<LoginDialogProps> = (props) => {
    const classes = useStyles();
    const [loggedIn, setLoggedIn] = React.useState(false);
    const [loginError, setLoginError] = React.useState("");
    const formik = useFormik({
        initialValues: {
            email: "",
            password: "",
        },
        onSubmit: values => {
            fetch(`${USER_SERVICE_URI}/auth`, {
                method: "POST",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(values),
            })
                .then(resp => {
                    if (resp.ok) {
                        setLoggedIn(true);
                    } else {
                        setLoggedIn(false);
                        resp.json().then(json => {
                            setLoginError(json.message);
                        }).catch(error => {
                            console.log(error);
                            setLoggedIn(false);
                            setLoginError("There was a problem logging in");
                        });
                    }
                })
                .catch(error => {
                    console.log(error);
                    setLoggedIn(false);
                    setLoginError("There was a problem logging in")
                });
        },
        validationSchema: Yup.object({
            email: Yup.string().email("Invalid email address").required("Requried"),
            password: Yup.string().min(6, "Must be at least 6 characterse"),
        }),
    });
    return (
        <Dialog {...props}>
            <DialogTitle>Login</DialogTitle>
            <form onSubmit={formik.handleSubmit} >
                <DialogContent>
                    <TextField
                        fullWidth
                        label="Email"
                        id="email"
                        name="email"
                        type="email"
                        onChange={formik.handleChange}
                        value={formik.values.email}
                        error={!!formik.errors.email}
                        helperText={formik.errors.email || null}
                    />

                    <TextField
                        fullWidth
                        label="Password"
                        id="password"
                        name="password"
                        type="password"
                        onChange={formik.handleChange}
                        value={formik.values.password}
                        error={!!formik.errors.password}
                        helperText={formik.errors.password || null}
                    />
                </DialogContent>
                <DialogActions>
                    <Container>
                        <div>
                            <Button variant="contained" type="submit" color="primary" className={classes.submitButton}>Submit</Button>
                        </div>
                        <div>
                            <Typography color="error">{loginError}</Typography>
                        </div>
                    </Container>
                </DialogActions>
            </form>
        </Dialog>
    )
}

interface SignUpDialogProps extends DialogProps { }

const SignUpDialog: FunctionComponent<SignUpDialogProps> = (props) => {
    console.log(props)
    const classes = useStyles();
    const [signedUp, setSignedUp] = React.useState(false);
    const [signUpError, setSignUpError] = React.useState("");
    const { handleSubmit, handleChange, values, errors } = useFormik({
        initialValues: {
            email: "",
            firstName: "",
            lastName: "",
            password: "",
        },
        onSubmit: values => {
            fetch(`${USER_SERVICE_URI}/auth`, {
                method: "POST",
                mode: "cors",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(values),
            })
                .then(resp => {
                    if (resp.ok) {
                        setSignedUp(true);
                    } else {
                        setSignedUp(false);
                        resp.json().then(json => {
                            setSignUpError(json.message);
                        }).catch(error => {
                            console.log(error);
                            setSignedUp(false);
                            setSignUpError("There was a problem signing up");
                        });
                    }
                })
                .catch(error => {
                    console.log(error);
                    setSignedUp(false);
                    setSignUpError("There was a problem signing up");
                });
        },
        validationSchema: Yup.object({
            email: Yup.string().email("Invalid email address").required("Requried"),
            firstName: Yup.string().required("Required"),
            lastName: Yup.string().required("Required"),
            password: Yup.string().min(6, "Must be at least 6 characterse").required("Required"),
        }),
    });
    return (
        <Dialog {...props} open={!signedUp && props.open}>
            <DialogTitle>Sign Up</DialogTitle>
            <form onSubmit={handleSubmit}>
                <DialogContent>
                    <TextField fullWidth label="Email" id="email" name="email" type="email"
                        onChange={handleChange} value={values.email} error={!!errors.email}
                        helperText={errors.email || null} />
                    <TextField fullWidth label="First name" id="firstName" name="firstName"
                        onChange={handleChange} value={values.firstName} error={!!errors.firstName}
                        helperText={errors.firstName || null} />
                    <TextField fullWidth label="Last name" id="lastName" name="lastName"
                        onChange={handleChange} value={values.lastName} error={!!errors.lastName}
                        helperText={errors.lastName || null} />

                    <TextField fullWidth label="Password" id="password" name="password" type="password"
                        onChange={handleChange} value={values.password} error={!!errors.password}
                        helperText={errors.password || null} />
                </DialogContent>

                <DialogActions>
                    <Grid container spacing={3}>
                        <Grid item xs>
                            <Button variant="contained" type="submit" color="primary" className={classes.submitButton}>
                                Submit
                            </Button>
                        </Grid>
                        <Grid item xs>
                            <Typography color="error">
                                {signUpError}
                            </Typography>
                        </Grid>
                    </Grid>
                </DialogActions>
            </form>
        </Dialog>
    )
}


const Home: FunctionComponent = () => {
    const classes = useStyles();
    const [openSignUp, setSignUp] = React.useState(false);
    const [openLogin, setLogin] = React.useState(false);

    const handleCloseSignUp = () => setSignUp(false);
    const handleToggleSignUp = () => setSignUp(!openSignUp);

    const handleCloseLogin = () => setLogin(false);
    const handleToggleLogin = () => setLogin(!openLogin);

    return (
        <div className={classes.root}>
            <AppBar position="static" color="primary">
                <Toolbar>
                    <Typography variant="h4" className={classes.title}></Typography>
                    <Button color="inherit" onClick={handleToggleSignUp}>Sign Up</Button>
                    <Button color="inherit" onClick={handleToggleLogin}>Login</Button>
                </Toolbar>
            </AppBar>
            <SignUpDialog className={classes.backdrop} open={openSignUp} onClose={handleCloseSignUp} />
            <LoginDialog onClose={handleCloseLogin} open={openLogin} className={classes.backdrop} />
            <Typography variant="h1" className={classes.title}>Badgerer</Typography>
            <Typography variant="subtitle1">
                Removing the work of managing badgework and attendance for scout Leaders.
            </Typography>
            <Box className={classes.blurb}>
                <Typography variant="body1">
                    Badgerer handles the book keeping around badgework and attendance. You'll no longer have to update your records for all scouts that,
                    attended. Instead, you just submit who attended and what badgework items were completed - badgerer handles the rest!
                </Typography>
            </Box>
        </div>
    );
}

export default App;
