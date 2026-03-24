import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export function CreateTaskDialog({ onCreate, disabled }: { onCreate: (title: string, description: string) => Promise<void>; disabled?: boolean }) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  const handleCreate = async () => {
    await onCreate(title, description);
    setTitle('');
    setDescription('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={disabled}>Nova tarefa</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Criar tarefa</DialogTitle>
          <DialogDescription>Essa ação já usa o endpoint real do backend.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="task-title">Título</Label>
            <Input id="task-title" value={title} onChange={(event) => setTitle(event.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="task-description">Descrição</Label>
            <Input id="task-description" value={description} onChange={(event) => setDescription(event.target.value)} />
          </div>
          <Button className="w-full" onClick={handleCreate}>Salvar</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
