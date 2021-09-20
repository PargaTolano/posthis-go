import React,{ useState, useEffect, useRef } from 'react';

import {
    Button,
    CssBaseline,
    TextField,
    Grid,
    Typography,
    Container,
    IconButton,
}from  '@material-ui/core';

import {
  Image as ImageIcon,
  AccountCircle
}from '@material-ui/icons'

import { makeStyles }     from '@material-ui/core/styles';

import { updateUser, validatePassword }                           from '_api';
import { fileToBase64, validateUpdateUser }     from '_utils';
import { handleResponse }                       from '_helpers';
import { authenticationService }                from '_services';    

import { UpdateUserViewModel }                  from '_model';

import defaultProfilePic  from 'assets/avatar-placeholder.svg';
import defaultCoverPic    from 'assets/background-placeholder.jpg';

const useStyles = makeStyles((theme) => ({
  paper: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.primary.main,
  },
  form: {
    width: '100%', 
    marginTop: theme.spacing(3),
  },
  submit: {
    margin: theme.spacing(3, 0, 2)
  },
  input:{
    display: 'none',
  },
  profilePicture:{
    width:            '100px',
    height:           '100px',
    objectFit:        'cover',
    borderRadius:     '50%',
    backgroundColor:  '#333333'
  },
  backgroundPicture:{
    width:            '100%',
    height:           '150px',
    objectFit:        'cover',
    borderRadius:     theme.spacing(0, 0, 2, 2),
    
    marginTop:theme.spacing(3),
  },
  pictures: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  userIcon:{
    color: '#ea5970',
    margin: theme.spacing(1),
  },
  inputImage:{
    display: 'none'
  },
  fieldWarning:{
    color: '#ea5970',
    fontSize: '0.8rem',
    marginTop: theme.spacing(1)
  },
  divider:{
    display: 'inline-block',
    width: '100%',
    height: 2,
    backgroundColor: '#DDDDDD',
    margin: theme.spacing(4, 0 ),
  }
}));

let timeout = null;

export const EditInfo = ( props ) => {

  const { user, setUser, handleClose } = props;

  const classes           = useStyles();

  const inputProfileRef   = useRef();
  const inputCoverRef     = useRef();

  const [state, setState] = useState({
    username:           user.username,
    tag:                user.tag,
    email:              user.email,
    password:           user.password,
    changedProfilePic:  false,
    profilePic:         null,
    profilePicPreview:  user.profilePicPath,
    changedCoverPic:    false,
    coverPic:           null,
    coverPicPreview:    user.coverPicPath,
    canSubmit:          false,
  });

  const [validation, setValidation]= useState({
    username:   false,
    tag:        false,
    email:      false,
    validated:  false,
  });

  useEffect(()=>{
    setValidation( validateUpdateUser( state ) );
  },[state]);

  const onChangeTextField = e =>{

    setState({
      ...state,
      [e.target.name]: e.target.value
    });

    if ( e.target.name == 'password'){

      if ( timeout != null ) clearTimeout(timeout);

      timeout = setTimeout(async (value)=> {

        if (value === '') return;

        const {data:responseData, err} = await validatePassword(value);

        if ( err !== null || responseData === null ) return;

        const { data } = responseData;

        setState({
          ...state,
          [e.target.name]: value,
          canSubmit: data
        });
        timeout = null;
      }, 300, e.target.value);
    }

  };

  const onChangeProfilePic = async (e) => {
    let file = e.target.files[0];

    if( !file )
      return;

    let url  = await fileToBase64(file);
    setState( x=>{
      let copy = {...x};
      copy.profilePicPreview = url;
      copy.changedProfilePic = true;
      copy.profilePic        = file;
      return copy;
    })
  };

  const onChangeCoverPic   = async (e) => {
    let file = e.target.files[0];

    if( !file )
      return;
      
    let url  = await fileToBase64(file);
    setState( x=>{
      let copy = {...x};
      copy.coverPicPreview = url;
      copy.changedCoverPic = true;
      copy.coverPic        = file;
      return copy;
    })
  };

  const onClickProfilePic = () => inputProfileRef.current?.click();
 
  const onClickCoverPic   = () => inputCoverRef  .current?.click();

  const onSubmit =async e => {
    e.preventDefault();

    let model = new UpdateUserViewModel(state);

    const {data:responseData, err } = await updateUser(authenticationService.currentUserValue.id, model);

    if ( err !== null || responseData === null ) return;

    const { data } = responseData;
    setUser( data );
    handleClose();
  };
  
  const submitable = validation.validated && state.canSubmit;

  return (
    <Container component='main' maxWidth='sm'>
      <CssBaseline />
      <div className={classes.paper}>

        <AccountCircle className={classes.userIcon}/>

        <Typography component='h1' variant='h5'>
          <strong>Your Profile</strong>
        </Typography>

        <Typography variant='body2'>
          Update Your Info!
        </Typography>

        <form className={classes.form} onSubmit={onSubmit} noValidate>

          <Grid container spacing={2}>

            <Grid item xs={12} className={classes.pictures}>
              <img className={classes.profilePicture} src={ state.profilePicPreview || defaultProfilePic}/>
            </Grid>
           
            <Grid item xs={12}   className={classes.pictures}>
              <input 
                accept='image/*' 
                className={classes.input} 
                type='file' ref={inputProfileRef} 
                onChange={onChangeProfilePic}
              />
              <label htmlFor='profile-button-file'> Foto de perfil
                <IconButton 
                 color='secondary' 
                 aria-label='upload picture' 
                 component='span' 
                 onClick={onClickProfilePic}
                >
                  <ImageIcon className={classes.imageIcon}/>
                </IconButton>
              </label>
            </Grid>
          
            <Grid item xs={12} sm={6}>
              <TextField
                name='username'
                autoComplete='fname'
                variant='outlined'
                required
                fullWidth
                value={state.username}
                onChange={onChangeTextField}
                label='Usuario'
                autoFocus
                error = {!validation.username}
                helperText = { !validation.username && 'valid username required' }
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                name='tag'
                variant='outlined'
                required
                fullWidth
                label='Tag'
                value={state.tag}
                onChange={onChangeTextField}
                autoComplete='tagname'
                error = {!validation.tag}
                helperText = { validation.tag ? '*just alphanumeric characters' : 'valid tag required' }
              />
            </Grid>

            <Grid item xs={12}>
              <TextField
                name='email'
                variant='outlined'
                required
                fullWidth
                label='Email'
                value={state.email}
                onChange = {onChangeTextField}
                autoComplete='email'
              />
              {
                !validation.email
                && 
                <Typography variant='body2' className={classes.fieldWarning}>
                * Email no valido
                </Typography>
              }
            </Grid>

            <Grid item xs={12} className={classes.pictures}>
              <img className={classes.backgroundPicture} src={ state.coverPicPreview || defaultCoverPic }/>
            </Grid>
    
            <Grid item xs={12} className={classes.pictures}>
              <input 
               accept='image/*' 
               className={classes.input} 
               type='file' 
               ref={inputCoverRef} 
               onChange={onChangeCoverPic}
              />
              <label htmlFor='background-button-file'> Foto de portada
                <IconButton 
                 color='secondary' 
                 aria-label='upload picture' 
                 component='span' 
                 onClick={onClickCoverPic}
                >
                  <ImageIcon className={classes.imageIcon}/>
                </IconButton>
              </label>
            </Grid>
            <Grid item xs={12} className={classes.pictures}>
              <div className={classes.divider}/>
              <TextField
                name='password'
                variant='outlined'
                required
                fullWidth
                label='Password'
                type='password'
                value={state.password}
                onChange = {onChangeTextField}
                autoComplete='email'
              />
            </Grid>

          </Grid>

          <Button
            type      ='submit'
            fullWidth
            variant   ='contained'
            color     ='primary'
            className ={classes.submit}
            disabled  = {!submitable}
          >
            Save
          </Button>
        </form>
      </div>
      
    </Container>
  );
}

export default  EditInfo;