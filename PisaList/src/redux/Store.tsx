import axios from 'axios'
interface Task {
    id: string;
    event: string;
    completed: boolean;
    is_cycle: boolean;
    description: string;
    importanceLevel:number;
    completed_Date: string;
}
interface Wish {
  id:string;
  event:string;
  is_cycle:boolean;
  description:string;
  is_shared:boolean;
}

// interface modifyTask{
//   id:string;
//   event:string;
//   description:string;
  
// }
import { configureStore } from '@reduxjs/toolkit';
import { createSlice } from '@reduxjs/toolkit';


interface TaskState {
  tasks: Task[];
}

interface WishState{
  wishes: Wish[];
}


const initialState: TaskState = {

    tasks: JSON.parse(localStorage.getItem('tasks') as string) || [],
    
};
  
const wishInitialState: WishState = {
    wishes: JSON.parse(localStorage.getItem('wishes') as string) || [],
};

const appearance=createSlice({
  name: 'appearance',
  initialState:{
    appear:false
  },
  reducers: {
    setAppearance:(state)=>{
      state.appear=!state.appear;
    }
  }
})

const modifyTask=createSlice({
  name: 'modifyTask',
  initialState: {
    modifyTask: {
      id:'',
      event: '',
      description: '',
      isCycle:false,
      completed:false,
      importanceLevel:0,
      completed_Date:''
    }
  },
  reducers: {
    // setTaskEvent: (state, action) => {
    //   state.modifyTask.event = action.payload;
    // },
    // setTaskDescription: (state, action) => {
    //   state.modifyTask.description = action.payload;
    // }
    modify:(state,action)=>{
      state.modifyTask=action.payload;
    }
  }
})
const taskSlice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {
    addTask: (state, action) => {
      state.tasks.push(action.payload);
      const token=localStorage.getItem('token');
      if (token){
        const add=async()=>{
          await axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/add`,
            JSON.stringify(action.payload),
            {headers:{Authorization:`Bearer ${token}`}},
          )
        };
        add();
      }else{
        localStorage.setItem('tasks', JSON.stringify(state.tasks));
      }
    },
    modify_Task:(state,action)=>{
      const index=state.tasks.findIndex(task=>task.id===action.payload.id);
      if(index!==-1){//判断是否找到
      state.tasks[index].event=action.payload.event;
      state.tasks[index].description=action.payload.description;
      state.tasks[index].is_cycle=action.payload.is_cycle;
      localStorage.setItem('tasks', JSON.stringify(state.tasks));
      }
    },
    deleteTask:(state, action)=>{
      const index=state.tasks.findIndex(task=>task.id===action.payload);
        state.tasks.splice(index, 1);
        const token=localStorage.getItem('token');
        if(token){
          const deleteTask=async()=>{
            await axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/delete`,
            JSON.stringify(action.payload),
            {headers:{Authorization:`Bearer ${token}`}},
          )
          };
          deleteTask();
        }else{
          localStorage.setItem('tasks', JSON.stringify(state.tasks));
        }
    },

    updateTask:(state, action)=>{//
        state.tasks=action.payload;
        const newTasks:Task[]=[]
        for (let i=0;i<state.tasks.length;i++){
          newTasks.push({
            id:state.tasks[i].id,
            event:state.tasks[i].event,
            description:state.tasks[i].description,            
            completed:state.tasks[i].completed,
            is_cycle:state.tasks[i].is_cycle,
            importanceLevel:i,//重新定义importanceLevel
            completed_Date:state.tasks[i].completed_Date,
          })
        }
        state.tasks=newTasks;
        const token=localStorage.getItem('token');
        if(token){
          const updateTask=async()=>{
            await axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/updateImportanceLevel`,
            JSON.stringify(state.tasks),
            {headers:{Authorization:`Bearer ${token}`}},
          )
          };
          updateTask();
        }else{
          localStorage.setItem('tasks', JSON.stringify(state.tasks));
        }
    },
    isCycle:(state, action)=>{
        // state.tasks[action.payload].is_cycle=!state.tasks[action.payload].is_cycle;
        const index=state.tasks.findIndex(task=>task.id===action.payload);
        state.tasks[index].is_cycle=!state.tasks[index].is_cycle;
        localStorage.setItem('tasks', JSON.stringify(state.tasks));

    },
    getTasks:(state,action)=>{
      state.tasks=action.payload;
    },

    finishTask:(state, action)=>{
      const date = new Date();
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      // state.tasks[action.payload].completed=true;
      // state.tasks[action.payload].completed_Date=`${year}-${month}-${day}`;
      const index=state.tasks.findIndex(task=>task.id===action.payload);
      state.tasks[index].completed=true;
      state.tasks[index].completed_Date=`${year}-${month}-${day}`;
      const token=localStorage.getItem('token');
      if(token){
        const finish=async()=>{
          await axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/update`,
            state.tasks,
            {headers:{Authorization:`Bearer ${token}`}},
          )
        };
        finish();
      }else{
      localStorage.setItem('tasks', JSON.stringify(state.tasks));
      }
    },
}
});

const wishSlice = createSlice({
  name: 'wishes',
  initialState: wishInitialState,
  reducers: {
    addWish: (state, action) => {
      state.wishes.push(action.payload);
      const token=localStorage.getItem('token');
      if(token){
        const add=async()=>{
          await axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/wish/add`,
            JSON.stringify(action.payload),
            {headers:{Authorization:`Bearer ${token}`}},
          )
        };
        add();
    }else{
      localStorage.setItem('wishes', JSON.stringify(state.wishes));
    }
  },
  getWishes:(state,action)=>{
    state.wishes=action.payload;
  },
  deleteWish:(state, action)=>{
    state.wishes.splice(action.payload, 1);
    const token=localStorage.getItem('token');
    if(token){
      const deleteWish=async()=>{
        axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/wish/delete`,
          action.payload,
          {headers:{Authorization:`Bearer ${token}`}},
        )
      };
      deleteWish();
    }else{
      localStorage.setItem('wishes', JSON.stringify(state.wishes));
    }
  },
  updateWishes:(state, action)=>{
    state.wishes=action.payload;
    const token=localStorage.getItem('token');
    if(token){
      const updateWish=async()=>{
        axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/wish/update`,
          action.payload,
          {headers:{Authorization:`Bearer ${token}`}},
        )
      };
      updateWish();
    }else{
      localStorage.setItem('wishes', JSON.stringify(state.wishes));
    }
  },
  switchCycle:(state, action)=>{
    state.wishes[action.payload].is_cycle=!state.wishes[action.payload].is_cycle;
    const token=localStorage.getItem('token');
    if(token){
      const switchCycle=async()=>{
        axios.post(`${import.meta.env.VITE_REACT_APP_BASE_URL}/todolist/wish/switchCycle`,
          action.payload,
          {headers:{Authorization:`Bearer ${token}`}},
        )
      };
      switchCycle();
    }else{
      localStorage.setItem('wishes', JSON.stringify(state.wishes));
    }
  },

}
});

const store = configureStore({
  reducer: {
    tasks: taskSlice.reducer,
    wishes: wishSlice.reducer,
    appearance: appearance.reducer,
    modifyTask: modifyTask.reducer,
  },
});

const { setAppearance } = appearance.actions;
const { modify } = modifyTask.actions;
const { addTask,modify_Task, deleteTask,updateTask,getTasks,finishTask,isCycle } = taskSlice.actions;
const {addWish,getWishes,deleteWish,updateWishes,switchCycle} =wishSlice.actions;
export { addTask,modify_Task, deleteTask,updateTask,getTasks,finishTask,isCycle,addWish,getWishes,deleteWish,updateWishes,setAppearance,switchCycle,modify };
export default store;


