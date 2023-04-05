import {
    Button,
    Center,
    Container,
    Flex,
    MantineProvider,
    Stack,
    Text,
    TextInput,
} from '@mantine/core'

import { useEffect, useState } from 'react'
import Todo from './components/Todo'
import { ITodo } from './Interfaces'

export default function App() {
    const [todos, setTodos] = useState<ITodo[]>([])
    const [title, setTitle] = useState<string>('')
    const [hasError, setHasError] = useState<boolean>(false)
    const handleAddTodo = async () => {
        setHasError(false)
        if (title.trim() === '') {
            setHasError(true)
            return
        }
        try {
            const request = await fetch('http://localhost:8080/todos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    title,
                }),
            })
            if (!request.ok) throw new Error('Something went wrong...')
            const data: ITodo = await request.json()
            const newTodo: ITodo = {
                ...data,
                hasError: false,
            }
            setTitle('')
            setTodos([...todos, newTodo])
        } catch (error) {
            alert(error)
        }
    }

    const handleEdit = async (id: string, newTitle: string) => {
        const hasError = newTitle.trim() === ''
        const newTodos = todos.map((todo) => {
            if (todo.id === id) {
                todo.title = newTitle
                todo.hasError = hasError
            }
            return todo
        })
        setTodos(newTodos)
        if (!hasError) {
            try {
                const request = await fetch(
                    `http://localhost:8080/todos/${id}`,
                    {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify({
                            title: newTitle,
                        }),
                    }
                )
                if (!request.ok) throw new Error('Something went wrong...')
            } catch (error) {
                alert(error)
            }
        }
    }

    const handleDelete = async (id: string) => {
        try {
            const request = await fetch(`http://localhost:8080/todos/${id}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
            })
            if (!request.ok) throw new Error('Something went wrong...')
            const newTodos = todos.filter((todo) => todo.id !== id)
            setTodos(newTodos)
        } catch (error) {
            alert(error)
        }
    }

    const handleComplete = async (id: string) => {
        try {
            const todo = todos.find((todo) => todo.id === id)
            if (!todo) throw new Error('Todo not found')

            const request = await fetch(`http://localhost:8080/todos/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    completed: !todo.completed,
                }),
            })

            if (!request.ok) throw new Error('Something went wrong...')
            const updatedTodos = todos.map((todo) => {
                if (todo.id === id) {
                    todo = { ...todo, completed: !todo.completed }
                }
                return todo
            })
            setTodos(updatedTodos)
        } catch (error) {
            alert(error)
        }
    }

    useEffect(() => {
        fetch('http://localhost:8080/todos')
            .then((res) => res.json())
            .then((data: ITodo[]) => {
                const final: ITodo[] = data.map((todo) => ({
                    ...todo,
                    hasError: false,
                }))
                setTodos(final)
            })
    }, [])
    return (
        <MantineProvider
            theme={{ colorScheme: 'dark' }}
            withGlobalStyles
            withNormalizeCSS
        >
            <Container size='xs'>
                <Center>
                    <Text size='xl' my='sm'>
                        Golang Todo App
                    </Text>
                </Center>
                <Flex align='flex-start' justify='center' mb='xl'>
                    <TextInput
                        error={hasError ? 'Title is required' : undefined}
                        placeholder='Walk the dog...'
                        onChange={(e) => setTitle(e.target.value)}
                        value={title}
                    />
                    <Button ml='sm' variant='outline' onClick={handleAddTodo}>
                        Add Todo
                    </Button>
                </Flex>

                <Stack>
                    {todos.map((todo) => (
                        <Todo
                            key={todo.id}
                            {...todo}
                            onEdit={handleEdit}
                            onDelete={handleDelete}
                            onComplete={handleComplete}
                        />
                    ))}
                </Stack>
            </Container>
        </MantineProvider>
    )
}
