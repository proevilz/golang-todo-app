import { Paper, Flex, Button, Text, TextInput } from '@mantine/core'
import { IconCheck, IconTrash } from '@tabler/icons-react'
import { ITodo } from '../Interfaces'
import { useEffect, useRef, useState } from 'react'
import { useClickOutside } from '@mantine/hooks'

interface IProps extends ITodo {
    onEdit: (id: string, title: string) => void
    onDelete: (id: string) => void
    onComplete: (id: string) => void
    onSave: (id: string) => void
}
const Todo = ({
    completed,
    title,
    id,
    hasError,
    onEdit,
    onDelete,
    onComplete,
    onSave,
}: IProps) => {
    const [editing, setEditing] = useState<boolean>(false)
    const [titleLength, setTitleLength] = useState<number>(0)
    const ref = useClickOutside(() => {
        if (editing && title.trim() !== '') {
            setEditing(false)
        }
    })
    useEffect(() => {
        if (editing) {
            setTitleLength(title.length)
        } else {
            if (title.length !== titleLength) {
                onSave(id)
            }
        }
    }, [editing])
    const inputRef = useRef<HTMLInputElement>(null)
    return (
        <Paper shadow='xs' p='sm' withBorder radius='sm' ref={ref}>
            <Flex justify='space-between' align='center'>
                <Paper
                    radius='xl'
                    withBorder
                    h='30px'
                    w='30px'
                    sx={() => ({ cursor: 'pointer' })}
                    onClick={() => onComplete(id)}
                >
                    <Flex justify='center' align='center' pt='2px'>
                        {completed ? <IconCheck color={'green'} /> : null}
                    </Flex>
                </Paper>
                <Text
                    sx={() => ({ flex: 1 })}
                    ml='sm'
                    strikethrough={completed}
                    onDoubleClick={() => setEditing(true)}
                >
                    <TextInput
                        ref={inputRef}
                        error={hasError}
                        value={title}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' && title.trim() !== '') {
                                setEditing(false)

                                inputRef?.current?.blur()
                            }
                        }}
                        variant={editing ? 'default' : 'unstyled'}
                        onChange={(e) => {
                            onEdit(id, e.target.value)
                        }}
                        sx={() => ({ pointerEvents: editing ? 'all' : 'none' })}
                    />
                </Text>

                <Button
                    compact
                    variant='subtle'
                    color='red'
                    onClick={() => onDelete(id)}
                >
                    <IconTrash />
                </Button>
            </Flex>
        </Paper>
    )
}

export default Todo
