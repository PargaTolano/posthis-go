import React,{ useEffect, useState } from 'react';

import { ReplyCard } from 'components/Reply';

import { replyService } from '_services';

import styles from '_styles/ReplyContainer.module.css';

export const ReplyContainer = ({id}) => {

    const [replies, setReplies] = useState();

    useEffect(()=>{ 
        const unsubscribe = replyService.reply$.subscribe( x => void setReplies(x));

        replyService.getPostReplies(id);

        return unsubscribe;
    }, []);

    return (
        <div className={styles.replyContainer}>
            {
                replies?.map((reply) => <ReplyCard  key={reply.replyID} id={id} reply={reply}/>)
            }
        </div>
    )
}

export default ReplyContainer;