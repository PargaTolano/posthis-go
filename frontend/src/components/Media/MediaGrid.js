import React from 'react';

import CancelIcon       from '@material-ui/icons/Cancel';

import styles from '_styles/CardGrid.module.css';

const MediaContainer = ( props ) => {

  const { media, state, setState } = props;

  const onClickDelete = ()=>{
    const index = state.medias.indexOf(media);
    if ( index === -1 )
      return;

    setState( x =>{
      let copy = {...x};
      copy.deleted = [...copy.deleted, copy.medias[index].mediaID];
      copy.medias = copy.medias.filter((elem, i)=> i !== index);
      return copy;
    });
  };

  return (
    <div className={styles.mediaContainer}>
      {
        media.isVideo ?
        (
          <video className={styles.mediaFit} controls>
              <source src={media.path} type={media.mime}/>
          </video>
        )
        :
        (<img src={media.path} className={styles.mediaFit}/>)
      }
      {
        state.editMode &&
          <CancelIcon className={styles.deleteIcon} onClick={onClickDelete}/>
      }
    </div>
  );
}


export const MediaGrid = ( props ) => {

  const { media, ...temp } = props;

  if ( media === null || media === undefined || media?.length === 0 )
    return null;

  const n = media?.length;

  

  return (
    <div className={styles.grid}>

        {
          media.map((v,i)=><MediaContainer key={v.id} media={v} {...temp}/>)
        }
    </div>
  );
};

export default MediaGrid;